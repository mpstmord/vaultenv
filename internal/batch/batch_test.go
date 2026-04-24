package batch_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/your-org/vaultenv/internal/batch"
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
		return nil, fmt.Errorf("path not found: %s", path)
	}
	return d, nil
}

func TestFetchAll_AllSuccess(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/a": {"key": "val-a"},
			"secret/b": {"key": "val-b"},
		},
	}
	b := batch.New(f, 2)
	reqs := []batch.Request{
		{Path: "secret/a", Field: "key"},
		{Path: "secret/b", Field: "key"},
	}
	results := b.FetchAll(context.Background(), reqs)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Value != "val-a" {
		t.Errorf("expected val-a, got %s", results[0].Value)
	}
	if results[1].Value != "val-b" {
		t.Errorf("expected val-b, got %s", results[1].Value)
	}
}

func TestFetchAll_PartialError(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/a": {"key": "val-a"},
		},
	}
	b := batch.New(f, 2)
	reqs := []batch.Request{
		{Path: "secret/a", Field: "key"},
		{Path: "secret/missing", Field: "key"},
	}
	results := b.FetchAll(context.Background(), reqs)
	if results[0].Err != nil {
		t.Errorf("expected no error for first result, got %v", results[0].Err)
	}
	if results[1].Err == nil {
		t.Error("expected error for missing path")
	}
}

func TestFetchAll_MissingField(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/a": {"other": "val"},
		},
	}
	b := batch.New(f, 1)
	reqs := []batch.Request{{Path: "secret/a", Field: "key"}}
	results := b.FetchAll(context.Background(), reqs)
	if results[0].Err == nil {
		t.Error("expected error for missing field")
	}
}

func TestFetchAll_UpstreamError(t *testing.T) {
	f := &mockFetcher{err: errors.New("vault unavailable")}
	b := batch.New(f, 2)
	reqs := []batch.Request{{Path: "secret/a", Field: "key"}}
	results := b.FetchAll(context.Background(), reqs)
	if results[0].Err == nil {
		t.Error("expected upstream error to propagate")
	}
}

func TestNew_DefaultWorkers(t *testing.T) {
	f := &mockFetcher{data: map[string]map[string]interface{}{}}
	b := batch.New(f, 0)
	if b == nil {
		t.Fatal("expected non-nil BatchFetcher")
	}
}

func TestNew_NilFetcherPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil upstream")
		}
	}()
	batch.New(nil, 2)
}
