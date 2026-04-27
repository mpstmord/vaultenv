package warmup_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/your-org/vaultenv/internal/warmup"
)

// --- fakes ---

type fakeFetcher struct {
	mu   sync.Mutex
	data map[string]map[string]interface{}
	err  error
}

func (f *fakeFetcher) GetSecretData(_ context.Context, path string) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return nil, f.err
	}
	return f.data[path], nil
}

type fakeWriter struct {
	mu   sync.Mutex
	vals map[string]string
}

func (w *fakeWriter) Set(key, value string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.vals == nil {
		w.vals = make(map[string]string)
	}
	w.vals[key] = value
}

// --- tests ---

func TestRun_Success(t *testing.T) {
	f := &fakeFetcher{
		data: map[string]map[string]interface{}{
			"secret/db": {"password": "s3cr3t"},
		},
	}
	w := &fakeWriter{}
	r := warmup.New(f, w)

	err := r.Run(context.Background(), []warmup.Mapping{
		{Path: "secret/db", Field: "password", EnvKey: "DB_PASSWORD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.vals["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", w.vals["DB_PASSWORD"])
	}
}

func TestRun_FetchError(t *testing.T) {
	f := &fakeFetcher{err: errors.New("vault unavailable")}
	w := &fakeWriter{}
	r := warmup.New(f, w)

	err := r.Run(context.Background(), []warmup.Mapping{
		{Path: "secret/db", Field: "password", EnvKey: "DB_PASSWORD"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_MissingField(t *testing.T) {
	f := &fakeFetcher{
		data: map[string]map[string]interface{}{
			"secret/db": {"user": "admin"},
		},
	}
	w := &fakeWriter{}
	r := warmup.New(f, w)

	err := r.Run(context.Background(), []warmup.Mapping{
		{Path: "secret/db", Field: "password", EnvKey: "DB_PASSWORD"},
	})
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestRun_NonStringField(t *testing.T) {
	f := &fakeFetcher{
		data: map[string]map[string]interface{}{
			"secret/cfg": {"port": 5432},
		},
	}
	w := &fakeWriter{}
	r := warmup.New(f, w)

	err := r.Run(context.Background(), []warmup.Mapping{
		{Path: "secret/cfg", Field: "port", EnvKey: "DB_PORT"},
	})
	if err == nil {
		t.Fatal("expected error for non-string field")
	}
}

func TestRun_EmptyMappings(t *testing.T) {
	f := &fakeFetcher{}
	w := &fakeWriter{}
	r := warmup.New(f, w)

	if err := r.Run(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NilFetcherPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil fetcher")
		}
	}()
	warmup.New(nil, &fakeWriter{})
}

func TestNew_NilWriterPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil writer")
		}
	}()
	warmup.New(&fakeFetcher{}, nil)
}
