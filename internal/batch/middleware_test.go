package batch_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/your-org/vaultenv/internal/batch"
)

func TestLoggingFetcher_LogsSuccess(t *testing.T) {
	f := &mockFetcher{
		data: map[string]map[string]interface{}{
			"secret/x": {"token": "abc"},
		},
	}
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	lf := batch.NewLoggingFetcher(f, logger)

	_, err := lf.GetSecretData(context.Background(), "secret/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "secret/x") {
		t.Errorf("expected log to contain path, got: %s", buf.String())
	}
}

func TestLoggingFetcher_LogsError(t *testing.T) {
	f := &mockFetcher{err: errors.New("boom")}
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	lf := batch.NewLoggingFetcher(f, logger)

	_, err := lf.GetSecretData(context.Background(), "secret/y")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(buf.String(), "error") {
		t.Errorf("expected error log, got: %s", buf.String())
	}
}

func TestLoggingFetcher_NilLoggerUsesDefault(t *testing.T) {
	f := &mockFetcher{data: map[string]map[string]interface{}{}}
	lf := batch.NewLoggingFetcher(f, nil)
	if lf == nil {
		t.Fatal("expected non-nil LoggingFetcher")
	}
}

func TestLoggingFetcher_NilUpstreamPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil upstream")
		}
	}()
	batch.NewLoggingFetcher(nil, nil)
}
