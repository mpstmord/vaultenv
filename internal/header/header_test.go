package header_test

import (
	"net/http"
	"testing"

	"github.com/your-org/vaultenv/internal/header"
)

func TestNew_ValidPairs(t *testing.T) {
	p, err := header.New([]string{"X-Request-ID: abc123", "X-Vault-Namespace: ns1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 2 {
		t.Fatalf("expected 2 headers, got %d", p.Len())
	}
}

func TestNew_EmptySlice(t *testing.T) {
	p, err := header.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 0 {
		t.Fatalf("expected 0 headers, got %d", p.Len())
	}
}

func TestNew_MalformedPair(t *testing.T) {
	_, err := header.New([]string{"NoColonHere"})
	if err == nil {
		t.Fatal("expected error for malformed pair, got nil")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := header.New([]string{": value"})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestNew_ValueWithColon(t *testing.T) {
	p, err := header.New([]string{"Authorization: Bearer tok:en"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h := p.Headers()
	if h["Authorization"] != "Bearer tok:en" {
		t.Fatalf("expected 'Bearer tok:en', got %q", h["Authorization"])
	}
}

func TestApply_SetsHeaders(t *testing.T) {
	p, err := header.New([]string{"X-Custom: hello", "X-Other: world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	p.Apply(req)
	if got := req.Header.Get("X-Custom"); got != "hello" {
		t.Errorf("X-Custom: expected 'hello', got %q", got)
	}
	if got := req.Header.Get("X-Other"); got != "world" {
		t.Errorf("X-Other: expected 'world', got %q", got)
	}
}

func TestApply_NilRequest(t *testing.T) {
	p, _ := header.New([]string{"X-Safe: yes"})
	// Must not panic.
	p.Apply(nil)
}

func TestHeaders_ReturnsCopy(t *testing.T) {
	p, _ := header.New([]string{"X-Foo: bar"})
	h := p.Headers()
	h["X-Foo"] = "mutated"
	if p.Headers()["X-Foo"] != "bar" {
		t.Error("Headers() should return an independent copy")
	}
}
