package quota_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/quota"
)

// stubFetcher is a minimal SecretFetcher for testing.
type stubFetcher struct {
	data map[string]interface{}
	err  error
	calls int
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	s.calls++
	return s.data, s.err
}

func TestGuardedFetcher_AllowsWithinLimit(t *testing.T) {
	stub := &stubFetcher{data: map[string]interface{}{"key": "val"}}
	l, _ := quota.New(3, time.Minute)
	gf := quota.NewGuardedFetcher(stub, l)

	for i := 0; i < 3; i++ {
		_, err := gf.GetSecretData(context.Background(), "secret/data/foo")
		if err != nil {
			t.Fatalf("request %d failed: %v", i+1, err)
		}
	}
	if stub.calls != 3 {
		t.Errorf("expected 3 upstream calls, got %d", stub.calls)
	}
}

func TestGuardedFetcher_BlocksOnExceeded(t *testing.T) {
	stub := &stubFetcher{data: map[string]interface{}{"k": "v"}}
	l, _ := quota.New(1, time.Minute)
	gf := quota.NewGuardedFetcher(stub, l)

	_, _ = gf.GetSecretData(context.Background(), "secret/data/bar")
	_, err := gf.GetSecretData(context.Background(), "secret/data/bar")
	if err == nil {
		t.Fatal("expected quota exceeded error")
	}
	var qe *quota.ErrQuotaExceeded
	if !errors.As(err, &qe) {
		t.Fatalf("expected *ErrQuotaExceeded, got %T", err)
	}
	// Upstream must not have been called the second time.
	if stub.calls != 1 {
		t.Errorf("expected 1 upstream call, got %d", stub.calls)
	}
}

func TestGuardedFetcher_PropagatesUpstreamError(t *testing.T) {
	upstreamErr := errors.New("vault unavailable")
	stub := &stubFetcher{err: upstreamErr}
	l, _ := quota.New(5, time.Minute)
	gf := quota.NewGuardedFetcher(stub, l)

	_, err := gf.GetSecretData(context.Background(), "secret/data/baz")
	if !errors.Is(err, upstreamErr) {
		t.Fatalf("expected upstream error, got %v", err)
	}
}

func TestGuardedFetcher_IndependentPaths(t *testing.T) {
	stub := &stubFetcher{data: map[string]interface{}{}}
	l, _ := quota.New(1, time.Minute)
	gf := quota.NewGuardedFetcher(stub, l)

	if _, err := gf.GetSecretData(context.Background(), "secret/data/p1"); err != nil {
		t.Fatalf("p1: %v", err)
	}
	if _, err := gf.GetSecretData(context.Background(), "secret/data/p2"); err != nil {
		t.Fatalf("p2: %v", err)
	}
}
