package services

import (
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
)

// SearchSession 搜索会话
type SearchSession struct {
	UserID      string              // 用户ID
	Keyword     string              // 搜索关键字
	Resources   []entity.Resource   // 搜索结果
	PageSize    int                 // 每页数量
	CurrentPage int                 // 当前页码
	TotalPages  int                 // 总页数
	LastAccess  time.Time           // 最后访问时间
}

// SearchSessionManager 搜索会话管理器
type SearchSessionManager struct {
	sessions map[string]*SearchSession  // 用户ID -> 搜索会话
	mutex    sync.RWMutex
}

// NewSearchSessionManager 创建搜索会话管理器
func NewSearchSessionManager() *SearchSessionManager {
	manager := &SearchSessionManager{
		sessions: make(map[string]*SearchSession),
	}

	// 启动清理过期会话的goroutine
	go manager.cleanupExpiredSessions()

	return manager
}

// CreateSession 创建或更新搜索会话
func (m *SearchSessionManager) CreateSession(userID, keyword string, resources []entity.Resource, pageSize int) *SearchSession {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session := &SearchSession{
		UserID:      userID,
		Keyword:     keyword,
		Resources:   resources,
		PageSize:    pageSize,
		CurrentPage: 1,
		TotalPages:  (len(resources) + pageSize - 1) / pageSize,
		LastAccess:  time.Now(),
	}

	m.sessions[userID] = session
	return session
}

// GetSession 获取搜索会话
func (m *SearchSessionManager) GetSession(userID string) *SearchSession {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[userID]
	if !exists {
		return nil
	}

	// 更新最后访问时间
	m.mutex.RUnlock()
	m.mutex.Lock()
	session.LastAccess = time.Now()
	m.mutex.Unlock()
	m.mutex.RLock()

	return session
}

// SetCurrentPage 设置当前页
func (m *SearchSessionManager) SetCurrentPage(userID string, page int) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[userID]
	if !exists {
		return false
	}

	if page < 1 || page > session.TotalPages {
		return false
	}

	session.CurrentPage = page
	session.LastAccess = time.Now()
	return true
}

// GetPageResources 获取指定页的资源
func (m *SearchSessionManager) GetPageResources(userID string, page int) []entity.Resource {
	m.mutex.RLock()
	session, exists := m.sessions[userID]
	m.mutex.RUnlock()

	if !exists {
		return nil
	}

	if page < 1 || page > session.TotalPages {
		return nil
	}

	start := (page - 1) * session.PageSize
	end := start + session.PageSize
	if end > len(session.Resources) {
		end = len(session.Resources)
	}

	// 更新当前页和最后访问时间
	m.mutex.Lock()
	session.CurrentPage = page
	session.LastAccess = time.Now()
	m.mutex.Unlock()

	return session.Resources[start:end]
}

// GetCurrentPageResources 获取当前页的资源
func (m *SearchSessionManager) GetCurrentPageResources(userID string) []entity.Resource {
	m.mutex.RLock()
	session, exists := m.sessions[userID]
	m.mutex.RUnlock()

	if !exists {
		return nil
	}

	return m.GetPageResources(userID, session.CurrentPage)
}

// HasNextPage 是否有下一页
func (m *SearchSessionManager) HasNextPage(userID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[userID]
	if !exists {
		return false
	}

	return session.CurrentPage < session.TotalPages
}

// HasPrevPage 是否有上一页
func (m *SearchSessionManager) HasPrevPage(userID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[userID]
	if !exists {
		return false
	}

	return session.CurrentPage > 1
}

// NextPage 下一页
func (m *SearchSessionManager) NextPage(userID string) []entity.Resource {
	m.mutex.Lock()
	session, exists := m.sessions[userID]
	m.mutex.Unlock()

	if !exists {
		return nil
	}

	if session.CurrentPage >= session.TotalPages {
		return nil
	}

	return m.GetPageResources(userID, session.CurrentPage+1)
}

// PrevPage 上一页
func (m *SearchSessionManager) PrevPage(userID string) []entity.Resource {
	m.mutex.Lock()
	session, exists := m.sessions[userID]
	m.mutex.Unlock()

	if !exists {
		return nil
	}

	if session.CurrentPage <= 1 {
		return nil
	}

	return m.GetPageResources(userID, session.CurrentPage-1)
}

// GetPageInfo 获取分页信息
func (m *SearchSessionManager) GetPageInfo(userID string) (currentPage, totalPages int, hasPrev, hasNext bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.sessions[userID]
	if !exists {
		return 0, 0, false, false
	}

	return session.CurrentPage, session.TotalPages, session.CurrentPage > 1, session.CurrentPage < session.TotalPages
}

// cleanupExpiredSessions 清理过期会话（超过1小时未访问）
func (m *SearchSessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		now := time.Now()
		for userID, session := range m.sessions {
			// 如果超过1小时未访问，清理该会话
			if now.Sub(session.LastAccess) > time.Hour {
				delete(m.sessions, userID)
			}
		}
		m.mutex.Unlock()
	}
}

// GlobalSearchSessionManager 全局搜索会话管理器
var GlobalSearchSessionManager = NewSearchSessionManager()