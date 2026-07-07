package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// TgSearchSession 一次 Telegram 搜索的分页状态（进程内、非持久化）。
// 用于在 Telegram callback_data 的 64 字节限制下，承载关键字/页大小等翻页所需上下文。
// 与 SearchSession（wechat 用）语义不同，故独立命名。
// Feature: 011-telegram-bot-enhance
type TgSearchSession struct {
	SessionID string
	ChatID    int64
	UserID    int64
	Keyword   string
	PageSize  int
	Total     int64
	CreatedAt time.Time
}

// TgSearchSessionStore Telegram 分页 session 存储抽象
type TgSearchSessionStore interface {
	Create(chatID, userID int64, keyword string, pageSize int, total int64) *TgSearchSession
	Get(sessionID string) (*TgSearchSession, bool)
	Delete(sessionID string)
}

// inMemoryTgSearchSessionStore 带 TTL + 容量上限（LRU）的进程内实现
type inMemoryTgSearchSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*TgSearchSession
	order    []string // 自旧至新，用于 LRU 淘汰
	maxSize  int
	ttl      time.Duration
}

// NewInMemoryTgSearchSessionStore 创建进程内分页 session 存储
func NewInMemoryTgSearchSessionStore(maxSize int, ttl time.Duration) TgSearchSessionStore {
	if maxSize <= 0 {
		maxSize = 1000
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &inMemoryTgSearchSessionStore{
		sessions: make(map[string]*TgSearchSession),
		maxSize:  maxSize,
		ttl:      ttl,
	}
}

// Create 创建新的搜索 session
func (s *inMemoryTgSearchSessionStore) Create(chatID, userID int64, keyword string, pageSize int, total int64) *TgSearchSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	sid := genSessionID()
	sess := &TgSearchSession{
		SessionID: sid,
		ChatID:    chatID,
		UserID:    userID,
		Keyword:   keyword,
		PageSize:  pageSize,
		Total:     total,
		CreatedAt: time.Now(),
	}
	s.sessions[sid] = sess
	s.order = append(s.order, sid)
	s.evictLocked()
	utils.Debug("[TELEGRAM:SESSION] 创建搜索会话 sid=%s keyword=%q total=%d (当前总数=%d)", sid, keyword, total, len(s.sessions))
	return sess
}

// Get 读取 session；过期则删除并返回未命中
func (s *inMemoryTgSearchSessionStore) Get(sessionID string) (*TgSearchSession, bool) {
	s.mu.RLock()
	sess, ok := s.sessions[sessionID]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Since(sess.CreatedAt) > s.ttl {
		s.Delete(sessionID)
		return nil, false
	}
	return sess, true
}

// Delete 删除指定 session
func (s *inMemoryTgSearchSessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
	for i, id := range s.order {
		if id == sessionID {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}

// evictLocked 清理过期与超容量条目（调用方持锁）
func (s *inMemoryTgSearchSessionStore) evictLocked() {
	now := time.Now()
	live := s.order[:0]
	for _, id := range s.order {
		sess, ok := s.sessions[id]
		if !ok {
			continue
		}
		if now.Sub(sess.CreatedAt) > s.ttl {
			delete(s.sessions, id)
			continue
		}
		live = append(live, id)
	}
	s.order = live
	// LRU：超容量则淘汰最旧
	for len(s.order) > s.maxSize {
		oldest := s.order[0]
		s.order = s.order[1:]
		delete(s.sessions, oldest)
	}
}

// genSessionID 生成 12 位十六进制随机 sessionID（48bit，碰撞概率极低）
func genSessionID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		utils.Error("[TELEGRAM:SESSION] 生成 sessionID 失败: %v", err)
		return hex.EncodeToString([]byte(time.Now().Format("060102150405")))[:12]
	}
	return hex.EncodeToString(b)
}
