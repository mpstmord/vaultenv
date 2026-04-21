package store

import (
	"context"
	"errors"
	"testing"
)

type mockFetcher struct {
	data map[string]map[string]interface{}
	err  error
}

func (m *mockFetcher) GetSecretData(_ context.Context, path string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	d, ok := m.data[path]
	if !ok {
		return nil, errors.New("not found")
	}
	return d, nil
}

func TestLoader_Load_Success(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/app": {"password": "s3cr3t"},
		},
	}
	s := New()
	l := NewLoader(f, s, []Mapping{
		{EnvKey: "APP_PASSWORD", Path: "secret/app", Field: "password"},
	})
	if err := l.Load(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := s.Get("APP_PASSWORD")
	if err != nil {
		t.Fatalf("key not stored: %v", err)
	}
	if v != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %s", v)
	}
}

func TestLoader_Load_FetchError(t *testing.T) {
	f := &mockFetcher{err: errors.New("vault unavailable")}
	s := New()
	l := NewLoader(f, s, []Mapping{
		{EnvKey: "X", Path: "secret/x", Field: "val"},
	})
	if err := l.Load(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoader_Load_MissingField(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/app": {"other": "value"},
		},
	}
	s := New()
	l := NewLoader(f, s, []Mapping{
		{EnvKey: "MISSING", Path: "secret/app", Field: "password"},
	})
	if err := l.Load(context.Background()); err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestLoader_Load_NonStringField(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/app": {"count": 42},
		},
	}
	s := New()
	l := NewLoader(f, s, []Mapping{
		{EnvKey: "COUNT", Path: "secret/app", Field: "count"},
	})
	if err := l.Load(context.Background()); err == nil {
		t.Fatal("expected error for non-string field")
	}
}

func TestLoader_Load_MultipleMappings(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/db": {"user": "admin", "pass": "hunter2"},
		},
	}
	s := New()
	l := NewLoader(f, s, []Mapping{
		{EnvKey: "DB_USER", Path: "secret/db", Field: "user"},
		{EnvKey: "DB_PASS", Path: "secret/db", Field: "pass"},
	})
	if err := l.Load(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", s.Len())
	}
}
