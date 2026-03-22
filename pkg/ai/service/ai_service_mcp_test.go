package service

import (
	"strings"
	"testing"
)

func TestValidateParamsSupportsRequiredStringSlice(t *testing.T) {
	as := &AIService{}
	tool := ToolDefinition{
		Name: "duckduckgo-search",
		Parameters: map[string]interface{}{
			"required": []string{"query"},
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	err := as.validateParams(tool, map[string]interface{}{})
	if err == nil || !strings.Contains(err.Error(), "缺少必需参数") {
		t.Fatalf("expected missing required param error, got: %v", err)
	}

	if err := as.validateParams(tool, map[string]interface{}{"query": "golang"}); err != nil {
		t.Fatalf("expected params to be valid, got error: %v", err)
	}
}

func TestParseToolCallsFromContentSupportsHyphenToolName(t *testing.T) {
	content := `<duckduckgo-search: {"query": "golang"}>`
	toolCalls := parseToolCallsFromContent(content, map[string]bool{
		"duckduckgo-search": true,
	})

	if len(toolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(toolCalls))
	}
	if toolCalls[0].Name != "duckduckgo-search" {
		t.Fatalf("expected tool name duckduckgo-search, got %s", toolCalls[0].Name)
	}
	if got, ok := toolCalls[0].Params["query"]; !ok || got != "golang" {
		t.Fatalf("expected query param golang, got: %v", toolCalls[0].Params)
	}
}

func TestCleanToolCallMarkersSupportsHyphenToolName(t *testing.T) {
	cleaned := cleanToolCallMarkers(`<duckduckgo-search/>继续回答`)
	if cleaned != "继续回答" {
		t.Fatalf("expected cleaned text to be 继续回答, got %q", cleaned)
	}
}
