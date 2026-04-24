package prefetch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/prefetch"
)

type countingFetcher struct {
	calls int
	data  map[string]interface{}
	err   error
}

func (c *countingFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	c.calls++
	if c.err != nil {
		return nil, c.err
	}
	return c.data, nil
}

func TestCachingFetcher_HitAvoidsFetch(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	cf := prefetch.NewCachingFetcher(upstream, time.Minute)

	_, _ = cf.GetSecretData(context.Background(), "secret/x")
	_, _ = cf.GetSecretData(context.Background(), "secret/x")

	if upstream.calls != 1 {
		t.Errorf("expected 1 upstream call, got %d", upstream.calls)
	}
}

func TestCachingFetcher_MissCallsUpstream(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	cf := prefetch.NewCachingFetcher(upstream, time.Minute)

	_, _ = cf.GetSecretData(context.Background(), "secret/a")
	_, _ = cf.GetSecretData(context.Background(), "secret/b")

	if upstream.calls != 2 {
		t.Errorf("expected 2 upstream calls, got %d", upstream.calls)
	}
}

func TestCachingFetcher_Expiry(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	cf := prefetch.NewCachingFetcher(upstream, time.Millisecond)

	_, _ = cf.GetSecretData(context.Background(), "secret/x")
	time.Sleep(5 * time.Millisecond)
	_, _ = cf.GetSecretData(context.Background(), "secret/x")

	if upstream.calls != 2 {
		t.Errorf("expected 2 upstream calls after expiry, got %d", upstream.calls)
	}
}

func TestCachingFetcher_Invalidate(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	cf := prefetch.NewCachingFetcher(upstream, time.Minute)

	_, _ = cf.GetSecretData(context.Background(), "secret/x")
	cf.Invalidate("secret/x")
	_, _ = cf.GetSecretData(context.Background(), "secret/x")

	if upstream.calls != 2 {
		t.Errorf("expected 2 upstream calls after invalidation, got %d", upstream.calls)
	}
}

func TestCachingFetcher_Flush(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	cf := prefetch.NewCachingFetcher(upstream, time.Minute)

	_, _ = cf.GetSecretData(context.Background(), "secret/a")
	_, _ = cf.GetSecretData(context.Background(), "secret/b")
	cf.Flush()
	_, _ = cf.GetSecretData(context.Background(), "secret/a")
	_, _ = cf.GetSecretData(context.Background(), "secret/b")

	if upstream.calls != 4 {
		t.Errorf("expected 4 upstream calls after flush, got %d", upstream.calls)
	}
}

func TestCachingFetcher_UpstreamError_NotCached(t *testing.T) {
	expected := errors.New("vault down")
	upstream := &countingFetcher{err: expected}
	cf := prefetch.NewCachingFetcher(upstream, time.Minute)

	_, err1 := cf.GetSecretData(context.Background(), "secret/x")
	_, err2 := cf.GetSecretData(context.Background(), "secret/x")

	if !errors.Is(err1, expected) || !errors.Is(err2, expected) {
		t.Error("expected upstream error to propagate")
	}
	if upstream.calls != 2 {
		t.Errorf("expected 2 calls when errors not cached, got %d", upstream.calls)
	}
}
