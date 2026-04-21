package audit

import (
	"errors"
	"testing"
	"time"
)

func TestNewEvent_Defaults(t *testing.T) {
	before := time.Now().UTC()
	ev := NewEvent(EventSecretFetch)
	after := time.Now().UTC()

	if ev.Type != EventSecretFetch {
		t.Errorf("expected type %q, got %q", EventSecretFetch, ev.Type)
	}
	if !ev.Success {
		t.Error("expected Success to be true by default")
	}
	if ev.Timestamp.Before(before) || ev.Timestamp.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ev.Timestamp, before, after)
	}
	if ev.Error != "" {
		t.Errorf("expected empty Error, got %q", ev.Error)
	}
}

func TestEvent_WithError(t *testing.T) {
	ev := NewEvent(EventError).WithError(errors.New("vault unreachable"))

	if ev.Success {
		t.Error("expected Success to be false after WithError")
	}
	if ev.Error != "vault unreachable" {
		t.Errorf("unexpected error string: %q", ev.Error)
	}
}

func TestEvent_WithNilError(t *testing.T) {
	ev := NewEvent(EventError).WithError(nil)

	if ev.Success {
		t.Error("expected Success false even with nil error")
	}
	if ev.Error != "" {
		t.Errorf("expected empty Error string for nil error, got %q", ev.Error)
	}
}

func TestEventType_Constants(t *testing.T) {
	types := []EventType{EventSecretFetch, EventSecretInject, EventProcessExec, EventError}
	seen := map[EventType]bool{}
	for _, et := range types {
		if seen[et] {
			t.Errorf("duplicate EventType value: %q", et)
		}
		seen[et] = true
	}
}

func TestEvent_WithErrorPreservesType(t *testing.T) {
	ev := NewEvent(EventSecretFetch).WithError(errors.New("permission denied"))

	if ev.Type != EventSecretFetch {
		t.Errorf("expected type %q to be preserved after WithError, got %q", EventSecretFetch, ev.Type)
	}
}

func TestEvent_WithErrorChaining(t *testing.T) {
	// Verify that calling WithError twice uses the most recent error.
	ev := NewEvent(EventSecretFetch).
		WithError(errors.New("first error")).
		WithError(errors.New("second error"))

	if ev.Error != "second error" {
		t.Errorf("expected last error %q, got %q", "second error", ev.Error)
	}
	if ev.Success {
		t.Error("expected Success to be false after chained WithError")
	}
}
