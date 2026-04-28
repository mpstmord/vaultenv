package resolve_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/resolve"
)

// stubFetcher is a test double for resolve.Fetcher.
type stubFetcher struct {
	data map[string]interface{}
	err  error
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return s.data, s.err
}

func TestNew_NoSources(t *testing.T) {
	_, err := resolve.New()
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNew_EmptySourceName(t *testing.T) {
	_, err := resolve.New(resolve.Source{Name: "", Fetcher: &stubFetcher{}})
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}

func TestNew_NilFetcher(t *testing.T) {
	_, err := resolve.New(resolve.Source{Name: "vault", Fetcher: nil})
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}

func TestResolve_FirstSourceSucceeds(t *testing.T) {
	want := map[string]interface{}{"key": "value"}
	r, err := resolve.New(
		resolve.Source{Name: "primary", Fetcher: &stubFetcher{data: want}},
		resolve.Source{Name: "fallback", Fetcher: &stubFetcher{data: map[string]interface{}{"key": "other"}}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, err := r.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != "primary" {
		t.Errorf("expected source %q, got %q", "primary", res.Source)
	}
	if res.Data["key"] != "value" {
		t.Errorf("unexpected data: %v", res.Data)
	}
}

func TestResolve_FallsBackOnError(t *testing.T) {
	want := map[string]interface{}{"token": "abc"}
	r, _ := resolve.New(
		resolve.Source{Name: "primary", Fetcher: &stubFetcher{err: errors.New("unavailable")}},
		resolve.Source{Name: "fallback", Fetcher: &stubFetcher{data: want}},
	)
	res, err := r.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != "fallback" {
		t.Errorf("expected source %q, got %q", "fallback", res.Source)
	}
}

func TestResolve_AllFail(t *testing.T) {
	r, _ := resolve.New(
		resolve.Source{Name: "a", Fetcher: &stubFetcher{err: errors.New("err-a")}},
		resolve.Source{Name: "b", Fetcher: &stubFetcher{err: errors.New("err-b")}},
	)
	_, err := r.Resolve(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error when all sources fail")
	}
}
