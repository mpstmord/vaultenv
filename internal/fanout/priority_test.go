package fanout_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/fanout"
)

func TestNewPriority_PanicsWithNoFetchers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic with no fetchers")
		}
	}()
	fanout.NewPriority()
}

func TestPriority_FirstSucceeds(t *testing.T) {
	a := &stubFetcher{data: map[string]interface{}{"src": "a"}}
	b := &stubFetcher{data: map[string]interface{}{"src": "b"}}
	p := fanout.NewPriority(a, b)

	got, err := p.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["src"] != "a" {
		t.Errorf("expected result from first fetcher, got %v", got["src"])
	}
}

func TestPriority_FallsBackOnError(t *testing.T) {
	bad := &stubFetcher{err: errors.New("unavailable")}
	good := &stubFetcher{data: map[string]interface{}{"src": "fallback"}}
	p := fanout.NewPriority(bad, good)

	got, err := p.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["src"] != "fallback" {
		t.Errorf("expected fallback result, got %v", got["src"])
	}
}

func TestPriority_AllFail_ReturnsLastError(t *testing.T) {
	a := &stubFetcher{err: errors.New("err-a")}
	b := &stubFetcher{err: errors.New("err-b")}
	p := fanout.NewPriority(a, b)

	_, err := p.GetSecretData(context.Background(), "secret/foo")
	if err == nil {
		t.Fatal("expected error when all fetchers fail")
	}
}

func TestPriority_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	bad := &stubFetcher{err: context.Canceled}
	good := &stubFetcher{data: map[string]interface{}{"k": "v"}}
	p := fanout.NewPriority(bad, good)

	_, err := p.GetSecretData(ctx, "secret/foo")
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}
