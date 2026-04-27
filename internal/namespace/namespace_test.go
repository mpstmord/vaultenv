package namespace_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/namespace"
)

// mockFetcher records the last path it received.
type mockFetcher struct {
	lastPath string
	data     map[string]interface{}
	err      error
}

func (m *mockFetcher) GetSecretData(_ context.Context, path string) (map[string]interface{}, error) {
	m.lastPath = path
	return m.data, m.err
}

func TestNew_EmptyPrefix(t *testing.T) {
	_, err := namespace.New("", &mockFetcher{})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNew_NilUpstreamPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream")
		}
	}()
	_, _ = namespace.New("prod", nil)
}

func TestNew_TrimsSlashes(t *testing.T) {
	n, err := namespace.New("/prod/", &mockFetcher{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := n.Prefix(); got != "prod" {
		t.Errorf("expected prefix \"prod\", got %q", got)
	}
}

func TestGetSecretData_PrependsPrefixToPath(t *testing.T) {
	mock := &mockFetcher{data: map[string]interface{}{"key": "val"}}
	n, _ := namespace.New("prod", mock)

	_, err := n.GetSecretData(context.Background(), "service/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.lastPath != "prod/service/db" {
		t.Errorf("expected \"prod/service/db\", got %q", mock.lastPath)
	}
}

func TestGetSecretData_StripsLeadingSlashFromPath(t *testing.T) {
	mock := &mockFetcher{data: map[string]interface{}{}}
	n, _ := namespace.New("staging", mock)

	_, _ = n.GetSecretData(context.Background(), "/secrets/token")
	if mock.lastPath != "staging/secrets/token" {
		t.Errorf("unexpected path: %q", mock.lastPath)
	}
}

func TestGetSecretData_PropagatesUpstreamError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := &mockFetcher{err: sentinel}
	n, _ := namespace.New("prod", mock)

	_, err := n.GetSecretData(context.Background(), "db/password")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestGetSecretData_ReturnsData(t *testing.T) {
	expected := map[string]interface{}{"username": "admin", "password": "s3cr3t"}
	mock := &mockFetcher{data: expected}
	n, _ := namespace.New("prod", mock)

	got, err := n.GetSecretData(context.Background(), "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["username"] != "admin" || got["password"] != "s3cr3t" {
		t.Errorf("unexpected data: %v", got)
	}
}
