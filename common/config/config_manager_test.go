package config

import (
	"os"
	"path/filepath"
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

func TestFileLoader_Load(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_rules.json")

	content := `{
		"version": 1,
		"updated_at": "2024-01-15 10:30:00",
		"description": "Test rules",
		"rules": [
			{
				"id": 1,
				"name": "Test Rule",
				"pan_id": 1,
				"domains": ["test.example.com"],
				"url_patterns": ["https?://test\\.example\\.com/s/([a-zA-Z0-9]+)"],
				"priority": 1,
				"enabled": true,
				"remark": "Test"
			}
		]
	}`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewFileLoader(testFile)

	config, err := loader.Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
		return
	}

	if config.Version != 1 {
		t.Errorf("Expected version 1, got %d", config.Version)
	}

	if len(config.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(config.Rules))
	}
}

func TestFileLoader_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.json")
	nonExistingFile := filepath.Join(tmpDir, "not_exists.json")

	err := os.WriteFile(existingFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader1 := NewFileLoader(existingFile)
	if !loader1.Exists() {
		t.Error("Expected file to exist")
	}

	loader2 := NewFileLoader(nonExistingFile)
	if loader2.Exists() {
		t.Error("Expected file to not exist")
	}
}

func TestConfigManager_LoadLocalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "pan_rules.json")

	content := `{
		"version": 1,
		"updated_at": "2024-01-15 10:30:00",
		"description": "Test config",
		"rules": [
			{
				"id": 1,
				"name": "Local Config Test",
				"pan_id": 2,
				"domains": ["local.example.com"],
				"url_patterns": ["https?://local\\.example\\.com/s/([a-zA-Z0-9]+)"],
				"priority": 1,
				"enabled": true,
				"remark": "Local test"
			}
		]
	}`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	manager := NewConfigManager()
	loader := NewFileLoader(testFile)
	manager.SetFileLoader(loader)

	err = manager.LoadLocalConfig()
	if err != nil {
		t.Errorf("LoadLocalConfig() error = %v", err)
		return
	}

	rule := manager.GetRuleByDomain("local.example.com")
	if rule == nil {
		t.Error("Expected rule from local config")
		return
	}

	if rule.PanID != 2 {
		t.Errorf("Expected PanID 2, got %d", rule.PanID)
	}
}

func TestConfigManager_LoadLocalConfig_FileNotExist(t *testing.T) {
	manager := NewConfigManager()
	loader := NewFileLoader("/nonexistent/path/rules.json")
	manager.SetFileLoader(loader)

	err := manager.LoadLocalConfig()
	if err != nil {
		t.Errorf("Expected no error when file not exists, got: %v", err)
	}
}