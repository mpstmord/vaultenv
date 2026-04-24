package prefetch_test

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/your-org/vaultenv/internal/prefetch"
)

type mockFetcher struct {
	data map[string]map[string]interface{}
	err  error
}

func (m *mockFetcher) GetSecretData(_ context.Context, path string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	if d, ok := m.data[path]; ok {
		return d, nil
	}
	return nil, errors.New("not found")
}

func TestNew_DefaultWorkers(t *testing.T) {
	f := &mockFetcher{}
	p := prefetch.New(f, 0)
	if p == nil {
		t.Fatal("expected non-nil prefetcher")
	}
}

func TestNew_NilFetcherPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fetcher")
		}
	}()
	prefetch.New(nil, 2)
}

func TestFetchAll_Success(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/a": {"key": "val-a"},
			"secret/b": {"key": "val-b"},
		},
	}
	p := prefetch.New(f, 2)
	results := p.FetchAll(context.Background(), []string{"secret/a", "secret/b"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	paths := []string{results[0].Path, results[1].Path}
	sort.Strings(paths)
	if paths[0] != "secret/a" || paths[1] != "secret/b" {
		t.Errorf("unexpected paths: %v", paths)
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.Path, r.Err)
		}
	}
}

func TestFetchAll_PartialError(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/ok": {"token": "abc"},
		},
	}
	p := prefetch.New(f, 2)
	results := p.FetchAll(context.Background(), []string{"secret/ok", "secret/missing"})
	var okResult, errResult prefetch.Result
	for _, r := range results {
		if r.Path == "secret/ok" {
			okResult = r
		} else {
			errResult = r
		}
	}
	if okResult.Err != nil {
		t.Errorf("expected no error for secret/ok, got %v", okResult.Err)
	}
	if errResult.Err == nil {
		t.Error("expected error for secret/missing")
	}
}

func TestFetchAll_EmptyPaths(t *testing.T) {
	f := &mockFetcher{}
	p := prefetch.New(f, 2)
	results := p.FetchAll(context.Background(), []string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFetchAll_FetcherError(t *testing.T) {
	expected := errors.New("vault unavailable")
	f := &mockFetcher{err: expected}
	p := prefetch.New(f, 2)
	results := p.FetchAll(context.Background(), []string{"secret/x"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !errors.Is(results[0].Err, expected) {
		t.Errorf("expected vault unavailable error, got %v", results[0].Err)
	}
}
