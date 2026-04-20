package watch

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestLogChangeHandler_Added(t *testing.T) {
	var buf bytes.Buffer
	h := LogChangeHandler(&buf)

	err := h(map[string]string{}, map[string]string{"TOKEN": "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "+ TOKEN (added)") {
		t.Fatalf("expected added entry, got: %s", buf.String())
	}
}

func TestLogChangeHandler_Removed(t *testing.T) {
	var buf bytes.Buffer
	h := LogChangeHandler(&buf)

	err := h(map[string]string{"OLD": "val"}, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "- OLD (removed)") {
		t.Fatalf("expected removed entry, got: %s", buf.String())
	}
}

func TestLogChangeHandler_Changed(t *testing.T) {
	var buf bytes.Buffer
	h := LogChangeHandler(&buf)

	err := h(map[string]string{"K": "v1"}, map[string]string{"K": "v2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "~ K (changed)") {
		t.Fatalf("expected changed entry, got: %s", buf.String())
	}
}

func TestLogChangeHandler_NoChange(t *testing.T) {
	var buf bytes.Buffer
	h := LogChangeHandler(&buf)

	err := h(map[string]string{"K": "v"}, map[string]string{"K": "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output for unchanged secrets, got: %s", buf.String())
	}
}

func TestChainChangeHandlers_StopsOnError(t *testing.T) {
	called := 0
	h1 := func(_, _ map[string]string) error { called++; return errors.New("boom") }
	h2 := func(_, _ map[string]string) error { called++; return nil }

	chain := ChainChangeHandlers(h1, h2)
	err := chain(nil, nil)

	if err == nil {
		t.Fatal("expected error from chain")
	}
	if called != 1 {
		t.Fatalf("expected chain to stop after first error, called=%d", called)
	}
}

func TestChainChangeHandlers_AllOK(t *testing.T) {
	called := 0
	h := func(_, _ map[string]string) error { called++; return nil }

	chain := ChainChangeHandlers(h, h, h)
	if err := chain(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 3 {
		t.Fatalf("expected all handlers called, got %d", called)
	}
}
