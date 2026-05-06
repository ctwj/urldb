package config

import (
	"testing"
)

func TestConfigManager_LoadHardcodedRules(t *testing.T) {
	manager := NewConfigManager()

	rules := []PanRuleConfig{
		{
			ID:          1,
			Name:        "Test Rule",
			PanID:       1,
			Domains:     []string{"test.example.com", "www.test.com"},
			URLPatterns: []string{`https?://test\.example\.com/s/([a-zA-Z0-9]+)`},
			Priority:    1,
			Enabled:     true,
		},
	}

	err := manager.LoadHardcodedRules(rules)
	if err != nil {
		t.Errorf("LoadHardcodedRules() error = %v", err)
	}

	if manager.GetConfigSource() != SourceHardcoded {
		t.Errorf("Expected config source to be SourceHardcoded")
	}

	if len(manager.GetAllRules()) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(manager.GetAllRules()))
	}
}

func TestConfigManager_GetRuleByDomain(t *testing.T) {
	manager := NewConfigManager()

	rules := []PanRuleConfig{
		{
			ID:          1,
			Name:        "Quark Rule",
			PanID:       1,
			Domains:     []string{"pan.quark.cn"},
			URLPatterns: []string{`https?://pan\.quark\.cn/s/([a-zA-Z0-9]+)`},
			Priority:    1,
			Enabled:     true,
		},
	}

	err := manager.LoadHardcodedRules(rules)
	if err != nil {
		t.Fatalf("LoadHardcodedRules() error = %v", err)
	}

	tests := []struct {
		name     string
		domain   string
		wantPanID int
	}{
		{"Exact match", "pan.quark.cn", 1},
		{"Subdomain match", "sub.pan.quark.cn", 1},
		{"No match", "example.com", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := manager.GetRuleByDomain(tt.domain)
			if tt.wantPanID == 0 {
				if rule != nil {
					t.Errorf("Expected nil rule for domain %s", tt.domain)
				}
			} else {
				if rule == nil {
					t.Errorf("Expected rule for domain %s, got nil", tt.domain)
					return
				}
				if rule.PanID != tt.wantPanID {
					t.Errorf("Expected PanID %d, got %d", tt.wantPanID, rule.PanID)
				}
			}
		})
	}
}

func TestConfigManager_Priority(t *testing.T) {
	manager := NewConfigManager()

	rules := []PanRuleConfig{
		{
			ID:          1,
			Name:        "Low Priority",
			PanID:       2,
			Domains:     []string{"priority.example.com"},
			URLPatterns: []string{`https?://priority\.example\.com/s/([a-zA-Z0-9]+)`},
			Priority:    10,
			Enabled:     true,
		},
		{
			ID:          2,
			Name:        "High Priority",
			PanID:       1,
			Domains:     []string{"priority.example.com"},
			URLPatterns: []string{`https?://priority\.example\.com/s/([a-zA-Z0-9]+)`},
			Priority:    1,
			Enabled:     true,
		},
	}

	err := manager.LoadHardcodedRules(rules)
	if err != nil {
		t.Fatalf("LoadHardcodedRules() error = %v", err)
	}

	rule := manager.GetRuleByDomain("priority.example.com")
	if rule == nil {
		t.Fatal("Expected rule, got nil")
	}

	if rule.PanID != 1 {
		t.Errorf("Expected high priority rule (PanID=1), got PanID=%d", rule.PanID)
	}
}

func TestConfigManager_DisabledRule(t *testing.T) {
	manager := NewConfigManager()

	rules := []PanRuleConfig{
		{
			ID:          1,
			Name:        "Disabled Rule",
			PanID:       1,
			Domains:     []string{"disabled.example.com"},
			URLPatterns: []string{`https?://disabled\.example\.com/s/([a-zA-Z0-9]+)`},
			Priority:    1,
			Enabled:     false,
		},
	}

	err := manager.LoadHardcodedRules(rules)
	if err != nil {
		t.Fatalf("LoadHardcodedRules() error = %v", err)
	}

	rule := manager.GetRuleByDomain("disabled.example.com")
	if rule != nil {
		t.Errorf("Expected nil for disabled rule, got %v", rule)
	}
}

func TestConfigManager_MultipleDomains(t *testing.T) {
	manager := NewConfigManager()

	rules := []PanRuleConfig{
		{
			ID:          1,
			Name:        "Multi-domain Rule",
			PanID:       1,
			Domains:     []string{"domain1.com", "domain2.com", "sub.domain3.com"},
			URLPatterns: []string{`https?://.*\.com/s/([a-zA-Z0-9]+)`},
			Priority:    1,
			Enabled:     true,
		},
	}

	err := manager.LoadHardcodedRules(rules)
	if err != nil {
		t.Fatalf("LoadHardcodedRules() error = %v", err)
	}

	tests := []struct {
		name     string
		domain   string
		wantPanID int
	}{
		{"Domain 1", "domain1.com", 1},
		{"Domain 2", "domain2.com", 1},
		{"Subdomain 3", "sub.domain3.com", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := manager.GetRuleByDomain(tt.domain)
			if rule == nil {
				t.Errorf("Expected rule for domain %s", tt.domain)
				return
			}
			if rule.PanID != tt.wantPanID {
				t.Errorf("Expected PanID %d, got %d", tt.wantPanID, rule.PanID)
			}
		})
	}
}