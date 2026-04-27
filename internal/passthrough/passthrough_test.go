package passthrough_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultenv/internal/passthrough"
)

// stubFetcher is a minimal Fetcher for testing.
type stubFetcher struct {
	data map[string]interface{}
	err  error
	calls int
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	s.calls++
	return s.data, s.err
}

func TestNew_PanicsOnNilUpstream(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream")
		}
	}()
	passthrough.New(nil, "")
}

func TestGetSecretData_DelegatesToUpstream(t *testing.T) {
	stub := &stubFetcher{data: map[string]interface{}{"key": "secret"}}
	pt := passthrough.New(stub, "VAULTENV_")

	got, err := pt.GetSecretData(context.Background(), "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "secret" {
		t.Errorf("expected 'secret', got %v", got["key"])
	}
	if stub.calls != 1 {
		t.Errorf("expected 1 upstream call, got %d", stub.calls)
	}
}

func TestGetSecretData_UsesEnvOverride(t *testing.T) {
	t.Setenv("VAULTENV_MYAPP_DB", "overridden")

	stub := &stubFetcher{data: map[string]interface{}{"key": "secret"}}
	pt := passthrough.New(stub, "VAULTENV_")

	got, err := pt.GetSecretData(context.Background(), "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["value"] != "overridden" {
		t.Errorf("expected 'overridden', got %v", got["value"])
	}
	if stub.calls != 0 {
		t.Errorf("upstream should not be called when env override is set")
	}
}

func TestGetSecretData_PropagatesUpstreamError(t *testing.T) {
	expected := errors.New("vault unavailable")
	stub := &stubFetcher{err: expected}
	pt := passthrough.New(stub, "")

	_, err := pt.GetSecretData(context.Background(), "some/path")
	if !errors.Is(err, expected) {
		t.Errorf("expected upstream error, got %v", err)
	}
}

func TestGetSecretData_NoPrefixEnvKey(t *testing.T) {
	t.Setenv("MYAPP_API_KEY", "no-prefix-value")

	stub := &stubFetcher{}
	pt := passthrough.New(stub, "")

	got, err := pt.GetSecretData(context.Background(), "myapp/api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["value"] != "no-prefix-value" {
		t.Errorf("expected 'no-prefix-value', got %v", got["value"])
	}
}
