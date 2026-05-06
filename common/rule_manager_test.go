package pan

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&PanRule{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestRuleManager_LoadRules(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ruleManager := NewRuleManager(db)

	err = ruleManager.LoadRules()
	if err != nil {
		t.Errorf("LoadRules() error = %v", err)
	}
}

func TestRuleManager_AddAndMatchRule(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Test Rule",
		PanID:       1,
		Domains:     "test.example.com,www.test.com",
		URLPatterns: `https?://test\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	tests := []struct {
		name     string
		url      string
		wantPan  int
		wantID   string
	}{
		{
			name:    "Exact domain match",
			url:     "https://test.example.com/s/abc123",
			wantPan: 1,
			wantID:  "abc123",
		},
		{
			name:    "Subdomain match",
			url:     "https://sub.test.example.com/s/def456",
			wantPan: 1,
			wantID:  "def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotPan, err := ruleManager.ExtractShareID(tt.url)
			if err != nil {
				t.Errorf("ExtractShareID() error = %v", err)
				return
			}
			if gotPan != tt.wantPan {
				t.Errorf("ExtractShareID() Pan = %v, want %v", gotPan, tt.wantPan)
			}
			if gotID != tt.wantID {
				t.Errorf("ExtractShareID() ID = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestRuleManager_Priority(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Low Priority Rule",
		PanID:       2,
		Domains:     "priority.example.com",
		URLPatterns: `https?://priority\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    10,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create low priority rule: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "High Priority Rule",
		PanID:       1,
		Domains:     "priority.example.com",
		URLPatterns: `https?://priority\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create high priority rule: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	gotID, gotPan, err := ruleManager.ExtractShareID("https://priority.example.com/s/test123")
	if err != nil {
		t.Errorf("ExtractShareID() error = %v", err)
		return
	}

	if gotPan != 1 {
		t.Errorf("Priority test: expected pan 1 (high priority), got %v", gotPan)
	}
	if gotID != "test123" {
		t.Errorf("Priority test: expected ID 'test123', got %v", gotID)
	}
}

func TestRuleManager_DisabledRule(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Disabled Rule",
		PanID:       1,
		Domains:     "disabled.example.com",
		URLPatterns: `https?://disabled\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     false,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create disabled rule: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	gotID, gotPan, err := ruleManager.ExtractShareID("https://disabled.example.com/s/test123")
	if err != nil {
		t.Errorf("ExtractShareID() error = %v", err)
		return
	}

	if gotPan != int(NotFound) {
		t.Errorf("Disabled rule should return NotFound, got %v", gotPan)
	}
	if gotID != "" {
		t.Errorf("Disabled rule should return empty ID, got %v", gotID)
	}
}

func TestRuleManager_MultipleDomains(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Multi-domain Rule",
		PanID:       1,
		Domains:     "domain1.com,domain2.com,sub.domain3.com",
		URLPatterns: `https?://.*\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create multi-domain rule: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	tests := []struct {
		name     string
		url      string
		wantPan  int
		wantID   string
	}{
		{"Domain 1", "https://domain1.com/s/aaa", 1, "aaa"},
		{"Domain 2", "https://domain2.com/s/bbb", 1, "bbb"},
		{"Subdomain 3", "https://sub.domain3.com/s/ccc", 1, "ccc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotPan, err := ruleManager.ExtractShareID(tt.url)
			if err != nil {
				t.Errorf("ExtractShareID() error = %v", err)
				return
			}
			if gotPan != tt.wantPan {
				t.Errorf("ExtractShareID() Pan = %v, want %v", gotPan, tt.wantPan)
			}
			if gotID != tt.wantID {
				t.Errorf("ExtractShareID() ID = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestRuleManager_ExtractServiceType(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Service Type Test",
		PanID:       2,
		Domains:     "servicetest.example.com",
		URLPatterns: `https?://servicetest\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create service type test rule: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	tests := []struct {
		name string
		url  string
		want ServiceType
	}{
		{"Match rule", "https://servicetest.example.com/s/test", ServiceType(2)},
		{"No match", "https://other.example.com/s/test", NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ruleManager.ExtractServiceTypeByRule(tt.url)
			if got != tt.want {
				t.Errorf("ExtractServiceTypeByRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleManager_Refresh(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ruleManager := NewRuleManager(db)
	err = ruleManager.LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	err = db.Create(&PanRule{
		Name:        "Refresh Test Rule",
		PanID:       1,
		Domains:     "refresh.example.com",
		URLPatterns: `https?://refresh\.example\.com/s/([a-zA-Z0-9]+)`,
		Priority:    1,
		Enabled:     true,
	}).Error
	if err != nil {
		t.Fatalf("Failed to create refresh test rule: %v", err)
	}

	err = ruleManager.Refresh()
	if err != nil {
		t.Errorf("Refresh() error = %v", err)
	}

	gotID, gotPan, err := ruleManager.ExtractShareID("https://refresh.example.com/s/test123")
	if err != nil {
		t.Errorf("ExtractShareID() after refresh error = %v", err)
		return
	}

	if gotPan != 1 {
		t.Errorf("After refresh: expected pan 1, got %v", gotPan)
	}
	if gotID != "test123" {
		t.Errorf("After refresh: expected ID 'test123', got %v", gotID)
	}
}