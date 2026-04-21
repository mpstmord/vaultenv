package drain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/drain"
)

// mockFetcher is a simple Fetcher stub.
type mockFetcher struct {
	data map[string]interface{}
	err  error
}

func (m *mockFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return m.data, m.err
}

func TestGuardedFetcher_AllowsBeforeDrain(t *testing.T) {
	upstream := &mockFetcher{data: map[string]interface{}{"key": "value"}}
	d := drain.New(time.Second)
	gf := drain.NewGuardedFetcher(upstream, d)

	data, err := gf.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("unexpected data: %v", data)
	}
}

func TestGuardedFetcher_BlocksAfterDrain(t *testing.T) {
	upstream := &mockFetcher{data: map[string]interface{}{"key": "value"}}
	d := drain.New(time.Second)
	gf := drain.NewGuardedFetcher(upstream, d)

	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("drain failed: %v", err)
	}

	_, err := gf.GetSecretData(context.Background(), "secret/foo")
	if !errors.Is(err, drain.ErrDraining) {
		t.Fatalf("expected ErrDraining, got: %v", err)
	}
}

func TestGuardedFetcher_PropagatesUpstreamError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	upstream := &mockFetcher{err: sentinel}
	d := drain.New(time.Second)
	gf := drain.NewGuardedFetcher(upstream, d)

	_, err := gf.GetSecretData(context.Background(), "secret/bar")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got: %v", err)
	}
}

func TestGuardedFetcher_HoldsAcquireDuringFetch(t *testing.T) {
	// The fetch blocks; Drain should wait for it to finish.
	blocked := make(chan struct{})
	unblock := make(chan struct{})

	upstream := &blockingFetcher{blocked: blocked, unblock: unblock}
	d := drain.New(2 * time.Second)
	gf := drain.NewGuardedFetcher(upstream, d)

	go func() { _, _ = gf.GetSecretData(context.Background(), "secret/x") }()

	// Wait until the fetch is in progress.
	<-blocked

	drainDone := make(chan error, 1)
	go func() { drainDone <- d.Drain(context.Background()) }()

	// Drain should not complete while fetch is blocked.
	select {
	case err := <-drainDone:
		t.Fatalf("Drain completed too early with: %v", err)
	case <-time.After(80 * time.Millisecond):
	}

	close(unblock) // let fetch finish

	if err := <-drainDone; err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}
}

type blockingFetcher struct {
	blocked chan struct{}
	unblock chan struct{}
}

func (b *blockingFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	close(b.blocked)
	<-b.unblock
	return nil, nil
}
