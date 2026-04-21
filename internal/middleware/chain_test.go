package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/middleware"
)

type stubFetcher struct {
	data map[string]interface{}
	err  error
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return s.data, s.err
}

func TestChain_NoMiddleware(t *testing.T) {
	base := &stubFetcher{data: map[string]interface{}{"k": "v"}}
	f := middleware.Chain(base)
	got, err := f.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected v, got %v", got["k"])
	}
}

func TestChain_OrderIsOutermostFirst(t *testing.T) {
	var order []int
	makeMiddleware := func(id int) middleware.Middleware {
		return func(next middleware.Fetcher) middleware.Fetcher {
			return middleware.FetcherFunc(func(ctx context.Context, path string) (map[string]interface{}, error) {
				order = append(order, id)
				return next.GetSecretData(ctx, path)
			})
		}
	}
	base := &stubFetcher{data: map[string]interface{}{}}
	f := middleware.Chain(base, makeMiddleware(1), makeMiddleware(2))
	_, _ = f.GetSecretData(context.Background(), "x")
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("unexpected order: %v", order)
	}
}

func TestChain_PropagatesError(t *testing.T) {
	expected := errors.New("boom")
	base := &stubFetcher{err: expected}
	f := middleware.Chain(base)
	_, err := f.GetSecretData(context.Background(), "x")
	if !errors.Is(err, expected) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}
