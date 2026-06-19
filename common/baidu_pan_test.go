package pan

import (
	"strings"
	"testing"
)

func TestUpdateCookieValue(t *testing.T) {
	tests := []struct {
		name   string
		cookie string
		key    string
		value  string
		// assertions on parsed result (order-independent, since map rebuild is unordered)
		want   string // value for key
		checks map[string]string
	}{
		{
			name:   "add new key to existing cookie",
			cookie: "BDUSS=abc; PANPSC=xxx",
			key:    "BDCLND",
			value:  "randsk123",
			want:   "randsk123",
			checks: map[string]string{"BDUSS": "abc", "PANPSC": "xxx", "BDCLND": "randsk123"},
		},
		{
			name:   "update existing key",
			cookie: "BDUSS=abc; BDCLND=old",
			key:    "BDCLND",
			value:  "new",
			want:   "new",
			checks: map[string]string{"BDUSS": "abc", "BDCLND": "new"},
		},
		{
			name:   "empty cookie",
			cookie: "",
			key:    "BDCLND",
			value:  "v",
			want:   "v",
			checks: map[string]string{"BDCLND": "v"},
		},
		{
			name:   "trim whitespace around pairs",
			cookie: " BDUSS = abc ;  STOKEN=tok ",
			key:    "BDCLND",
			value:  "v",
			want:   "v",
			checks: map[string]string{"BDUSS": "abc", "STOKEN": "tok", "BDCLND": "v"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateCookieValue(tt.cookie, tt.key, tt.value)
			// parse result back into a map (order-independent)
			m := map[string]string{}
			for _, pair := range strings.Split(got, ";") {
				pair = strings.TrimSpace(pair)
				if pair == "" {
					continue
				}
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) != 2 {
					t.Fatalf("malformed pair %q in result %q", pair, got)
				}
				m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
			if m[tt.key] != tt.want {
				t.Fatalf("key %s = %q, want %q (result=%q)", tt.key, m[tt.key], tt.want, got)
			}
			for k, v := range tt.checks {
				if m[k] != v {
					t.Fatalf("key %s = %q, want %q (result=%q)", k, m[k], v, got)
				}
			}
		})
	}
}
