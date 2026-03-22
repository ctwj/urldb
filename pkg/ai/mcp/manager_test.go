package mcp

import "testing"

func TestParseConfigContentValidatesWithoutSideEffects(t *testing.T) {
	content := `{
		"mcpServers": {
			"demo": {
				"transport": "stdio",
				"command": "echo",
				"args": ["ok"],
				"enabled": true,
				"auto_start": false
			}
		}
	}`

	cfg, err := parseConfigContent(content)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	if _, exists := cfg.MCPServers["demo"]; !exists {
		t.Fatalf("expected service demo in parsed config")
	}
}

func TestReloadConfigDoesNotRegisterMockTools(t *testing.T) {
	manager := NewMCPManager()
	content := `{
		"mcpServers": {
			"demo": {
				"transport": "stdio",
				"command": "echo",
				"args": ["ok"],
				"enabled": true,
				"auto_start": false
			}
		}
	}`

	if err := manager.ReloadConfig(content); err != nil {
		t.Fatalf("reload config failed: %v", err)
	}

	tools := manager.GetToolRegistry().GetTools("demo")
	if len(tools) != 0 {
		t.Fatalf("expected no tools to be registered without start, got %d", len(tools))
	}

	statuses := manager.ListServiceStatuses()
	status, exists := statuses["demo"]
	if !exists {
		t.Fatalf("expected service status demo to exist")
	}
	if status.Status != "stopped" {
		t.Fatalf("expected service status stopped, got %s", status.Status)
	}
}

func TestStopClientIsIdempotentForStoppedService(t *testing.T) {
	manager := NewMCPManager()
	manager.services["demo"] = &ServiceStatus{
		Name:   "demo",
		Status: "stopped",
		Tools:  []Tool{},
	}

	if err := manager.StopClient("demo"); err != nil {
		t.Fatalf("expected stop to be idempotent, got error: %v", err)
	}
	if manager.services["demo"].Status != "stopped" {
		t.Fatalf("expected service to remain stopped, got %s", manager.services["demo"].Status)
	}
}
