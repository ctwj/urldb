package pan

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/ctwj/urldb/db/entity"
)

// PanRuleInfo 运行时规则信息
type PanRuleInfo struct {
	ID          uint
	PanID       uint
	Name        string
	Domains     []string
	URLPatterns []*regexp.Regexp
	Priority    int
	Enabled     bool
	PanKey      int
	PanName     string
}

// RuleManager 规则管理器
type RuleManager struct {
	rules      map[string]*PanRuleInfo
	domainMap  map[string]*PanRuleInfo
	mu         sync.RWMutex
	db         *gorm.DB
	stopChan   chan struct{}
}

// NewRuleManager 创建规则管理器
func NewRuleManager(db *gorm.DB) *RuleManager {
	return &RuleManager{
		rules:     make(map[string]*PanRuleInfo),
		domainMap: make(map[string]*PanRuleInfo),
		db:        db,
		stopChan:  make(chan struct{}),
	}
}

// LoadRules 从数据库加载规则
func (m *RuleManager) LoadRules() error {
	var panRules []entity.PanRule
	err := m.db.Preload("Pan").Where("enabled = ?", true).Order("priority ASC").Find(&panRules).Error
	if err != nil {
		return err
	}

	newRules := make(map[string]*PanRuleInfo)
	newDomainMap := make(map[string]*PanRuleInfo)

	for _, rule := range panRules {
		domains := strings.Split(rule.Domains, ",")
		for i := range domains {
			domains[i] = strings.TrimSpace(domains[i])
		}

		patternStrs := strings.Split(rule.URLPatterns, ",")
		var patterns []*regexp.Regexp
		for _, patternStr := range patternStrs {
			patternStr = strings.TrimSpace(patternStr)
			if patternStr != "" {
				pattern, err := regexp.Compile(patternStr)
				if err == nil {
					patterns = append(patterns, pattern)
				}
			}
		}

		ruleInfo := &PanRuleInfo{
			ID:          rule.ID,
			PanID:       rule.PanID,
			Name:        rule.Name,
			Domains:     domains,
			URLPatterns: patterns,
			Priority:    rule.Priority,
			Enabled:     rule.Enabled,
			PanKey:      rule.Pan.Key,
			PanName:     rule.Pan.Name,
		}

		newRules[fmt.Sprintf("%d", rule.ID)] = ruleInfo

		for _, domain := range domains {
			if existing, ok := newDomainMap[domain]; ok {
				if rule.Priority < existing.Priority {
					newDomainMap[domain] = ruleInfo
				}
			} else {
				newDomainMap[domain] = ruleInfo
			}
		}
	}

	m.mu.Lock()
	m.rules = newRules
	m.domainMap = newDomainMap
	m.mu.Unlock()

	return nil
}

// GetRuleByDomain 根据域名获取规则
func (m *RuleManager) GetRuleByDomain(domain string) *PanRuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	domain = strings.ToLower(strings.TrimSpace(domain))

	if rule, ok := m.domainMap[domain]; ok {
		return rule
	}

	for d, rule := range m.domainMap {
		if strings.HasSuffix(domain, d) || strings.HasSuffix(d, domain) {
			return rule
		}
	}

	return nil
}

// GetRuleByPanKey 根据网盘Key获取规则列表
func (m *RuleManager) GetRuleByPanKey(panKey int) []*PanRuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PanRuleInfo
	for _, rule := range m.rules {
		if rule.PanKey == panKey && rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// GetAllRules 获取所有规则
func (m *RuleManager) GetAllRules() []*PanRuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PanRuleInfo
	for _, rule := range m.rules {
		result = append(result, rule)
	}
	return result
}

// Refresh 手动刷新规则
func (m *RuleManager) Refresh() error {
	return m.LoadRules()
}

// StartAutoRefresh 启动自动刷新
func (m *RuleManager) StartAutoRefresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.LoadRules()
			case <-m.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// StopAutoRefresh 停止自动刷新
func (m *RuleManager) StopAutoRefresh() {
	close(m.stopChan)
}

// ExtractShareID 从URL中提取分享ID
func (m *RuleManager) ExtractShareID(urlStr string) (string, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", int(NotFound), err
	}

	domain := parsedURL.Hostname()
	domain = strings.ToLower(domain)

	for _, rule := range m.rules {
		if !rule.Enabled {
			continue
		}

		for _, d := range rule.Domains {
			if strings.EqualFold(d, domain) || strings.HasSuffix(domain, d) {
				for _, pattern := range rule.URLPatterns {
					matches := pattern.FindStringSubmatch(urlStr)
					if len(matches) >= 2 {
						return matches[1], rule.PanKey, nil
					}
				}
			}
		}
	}

	return "", int(NotFound), nil
}

// ExtractServiceTypeByRule 根据规则识别服务类型
func (m *RuleManager) ExtractServiceTypeByRule(urlStr string) ServiceType {
	m.mu.RLock()
	defer m.mu.RUnlock()

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return NotFound
	}

	domain := strings.ToLower(parsedURL.Hostname())

	for _, rule := range m.rules {
		if !rule.Enabled {
			continue
		}

		for _, d := range rule.Domains {
			if strings.EqualFold(d, domain) || strings.HasSuffix(domain, d) {
				return ServiceType(rule.PanKey)
			}
		}
	}

	return NotFound
}