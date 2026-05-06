package config

import (
	"regexp"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
)

// RuleInfo 运行时规则信息
type RuleInfo struct {
	ID          int
	PanID       int
	Name        string
	Domains     []string
	URLPatterns []*regexp.Regexp
	Priority    int
	Enabled     bool
	Remark      string
}

// ConfigManager 配置管理器
type ConfigManager struct {
	rules        map[int]*RuleInfo
	domainMap    map[string]*RuleInfo
	source       ConfigSource
	remoteFetcher *RemoteFetcher
	mu           sync.RWMutex
	stopChan     chan struct{}
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		rules:      make(map[int]*RuleInfo),
		domainMap:  make(map[string]*RuleInfo),
		source:     SourceHardcoded,
		stopChan:   make(chan struct{}),
	}
}

// SetRemoteFetcher 设置远程配置获取器
func (m *ConfigManager) SetRemoteFetcher(fetcher *RemoteFetcher) {
	m.remoteFetcher = fetcher
}

// LoadHardcodedRules 加载硬编码规则
func (m *ConfigManager) LoadHardcodedRules(rules []PanRuleConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rules = make(map[int]*RuleInfo)
	m.domainMap = make(map[string]*RuleInfo)

	for _, ruleConfig := range rules {
		rule, err := m.parseRule(ruleConfig)
		if err != nil {
			utils.Warn("Failed to parse hardcoded rule '%s': %v", ruleConfig.Name, err)
			continue
		}
		m.addRule(rule)
	}

	m.source = SourceHardcoded
	utils.Info("Loaded %d hardcoded rules", len(m.rules))

	return nil
}

// LoadRemoteConfig 加载远程配置
func (m *ConfigManager) LoadRemoteConfig() error {
	if m.remoteFetcher == nil {
		return nil
	}

	config, err := m.remoteFetcher.Fetch()
	if err != nil {
		utils.Warn("Failed to fetch remote config: %v", err)
		return err
	}

	return m.ApplyRemoteConfig(config)
}

// ApplyRemoteConfig 应用远程配置
func (m *ConfigManager) ApplyRemoteConfig(config *RemoteConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	newRules := make(map[int]*RuleInfo)
	newDomainMap := make(map[string]*RuleInfo)

	for _, ruleConfig := range config.Rules {
		rule, err := m.parseRule(PanRuleConfig(ruleConfig))
		if err != nil {
			utils.Warn("Failed to parse remote rule '%s': %v", ruleConfig.Name, err)
			continue
		}
		newRules[rule.ID] = rule

		for _, domain := range rule.Domains {
			if existing, ok := newDomainMap[domain]; ok {
				if rule.Priority < existing.Priority {
					newDomainMap[domain] = rule
				}
			} else {
				newDomainMap[domain] = rule
			}
		}
	}

	m.rules = newRules
	m.domainMap = newDomainMap
	m.source = SourceRemote

	utils.Info("Applied remote config version %d, loaded %d rules", config.Version, len(m.rules))

	return nil
}

// TryUpdateFromRemote 尝试从远程更新配置
func (m *ConfigManager) TryUpdateFromRemote() (bool, error) {
	if m.remoteFetcher == nil {
		return false, nil
	}

	hasUpdate, newVersion, err := m.remoteFetcher.CheckUpdate()
	if err != nil {
		utils.Warn("Failed to check update: %v", err)
		return false, err
	}

	if hasUpdate {
		utils.Info("Found new config version: %d (current: %d)", newVersion, m.remoteFetcher.GetLastVersion())
		err := m.LoadRemoteConfig()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// StartAutoUpdate 启动自动更新检查
func (m *ConfigManager) StartAutoUpdate(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.TryUpdateFromRemote()
			case <-m.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// StopAutoUpdate 停止自动更新检查
func (m *ConfigManager) StopAutoUpdate() {
	close(m.stopChan)
}

// GetRuleByDomain 根据域名获取规则
func (m *ConfigManager) GetRuleByDomain(domain string) *RuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if rule, ok := m.domainMap[domain]; ok && rule.Enabled {
		return rule
	}

	for d, rule := range m.domainMap {
		if len(domain) >= len(d) && domain[len(domain)-len(d):] == d && rule.Enabled {
			return rule
		}
	}

	return nil
}

// GetRuleByPanID 根据网盘ID获取规则列表
func (m *ConfigManager) GetRuleByPanID(panID int) []*RuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*RuleInfo
	for _, rule := range m.rules {
		if rule.PanID == panID && rule.Enabled {
			result = append(result, rule)
		}
	}
	return result
}

// GetAllRules 获取所有规则
func (m *ConfigManager) GetAllRules() []*RuleInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*RuleInfo
	for _, rule := range m.rules {
		result = append(result, rule)
	}
	return result
}

// GetConfigSource 获取当前配置来源
func (m *ConfigManager) GetConfigSource() ConfigSource {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.source
}

// GetRemoteVersion 获取远程配置版本
func (m *ConfigManager) GetRemoteVersion() int {
	if m.remoteFetcher == nil {
		return 0
	}
	return m.remoteFetcher.GetLastVersion()
}

// parseRule 解析规则配置
func (m *ConfigManager) parseRule(config PanRuleConfig) (*RuleInfo, error) {
	var patterns []*regexp.Regexp
	for _, patternStr := range config.URLPatterns {
		pattern, err := regexp.Compile(patternStr)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}

	return &RuleInfo{
		ID:          config.ID,
		PanID:       config.PanID,
		Name:        config.Name,
		Domains:     config.Domains,
		URLPatterns: patterns,
		Priority:    config.Priority,
		Enabled:     config.Enabled,
		Remark:      config.Remark,
	}, nil
}

// addRule 添加规则到映射
func (m *ConfigManager) addRule(rule *RuleInfo) {
	m.rules[rule.ID] = rule

	for _, domain := range rule.Domains {
		if existing, ok := m.domainMap[domain]; ok {
			if rule.Priority < existing.Priority {
				m.domainMap[domain] = rule
			}
		} else {
			m.domainMap[domain] = rule
		}
	}
}