package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultenv/internal/audit"
)

func TestLog_Disabled(t *testing.T) {
	l := audit.NewLogger(nil)
	if err := l.Log(audit.Event{Type: audit.EventSecretFetch, Success: true}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.Log(audit.Event{
		Type:    audit.EventSecretFetch,
		Path:    "secret/data/myapp",
		Success: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Event
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Type != audit.EventSecretFetch {
		t.Errorf("expected type %q, got %q", audit.EventSecretFetch, got.Type)
	}
	if got.Path != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", got.Path)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_ErrorEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	_ = l.Log(audit.Event{
		Type:    audit.EventSecretResolve,
		Path:    "secret/data/myapp",
		Field:   "password",
		Success: false,
		Error:   "field not found",
	})

	var got audit.Event
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Error != "field not found" {
		t.Errorf("unexpected error field: %q", got.Error)
	}
}

func TestNilLogger_Safe(t *testing.T) {
	var l *audit.Logger
	if err := l.Log(audit.Event{Type: audit.EventExec}); err != nil {
		t.Fatalf("nil logger should be safe, got %v", err)
	}
}
