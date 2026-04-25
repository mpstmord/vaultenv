package fanout_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/fanout"
)

type stubFetcher struct {
	data map[string]interface{}
	err  error
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return s.data, s.err
}

func TestNew_PanicsWithNoFetchers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic with no fetchers")
		}
	}()
	fanout.New()
}

func TestGetSecretData_SingleFetcher(t *testing.T) {
	f := fanout.New(&stubFetcher{data: map[string]interface{}{"key": "val"}})
	got, err := f.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "val" {
		t.Errorf("expected val, got %v", got["key"])
	}
}

func TestGetSecretData_MergesMultipleFetchers(t *testing.T) {
	a := &stubFetcher{data: map[string]interface{}{"a": "1", "shared": "from-a"}}
	b := &stubFetcher{data: map[string]interface{}{"b": "2", "shared": "from-b"}}
	f := fanout.New(a, b)

	got, err := f.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["a"] != "1" || got["b"] != "2" {
		t.Errorf("missing keys: %v", got)
	}
	if got["shared"] == nil {
		t.Error("shared key should be present")
	}
}

func TestGetSecretData_ErrorPropagates(t *testing.T) {
	ok := &stubFetcher{data: map[string]interface{}{"k": "v"}}
	bad := &stubFetcher{err: errors.New("vault unreachable")}
	f := fanout.New(ok, bad)

	_, err := f.GetSecretData(context.Background(), "secret/foo")
	if err == nil {
		t.Fatal("expected error from failing fetcher")
	}
}

func TestGetSecretData_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	f := fanout.New(&stubFetcher{err: context.Canceled})
	_, err := f.GetSecretData(ctx, "secret/foo")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
