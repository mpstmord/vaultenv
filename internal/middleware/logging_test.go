package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/your-org/vaultenv/internal/middleware"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestLoggingMiddleware_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	base := &stubFetcher{data: map[string]interface{}{"user": "admin"}}
	f := middleware.Chain(base, middleware.NewLoggingMiddleware(logger))

	_, err := f.GetSecretData(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "secret fetched") {
		t.Errorf("expected log line 'secret fetched', got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "secret/app") {
		t.Errorf("expected path in log, got: %s", buf.String())
	}
}

func TestLoggingMiddleware_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	base := &stubFetcher{err: errors.New("vault unavailable")}
	f := middleware.Chain(base, middleware.NewLoggingMiddleware(logger))

	_, err := f.GetSecretData(context.Background(), "secret/db")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(buf.String(), "secret fetch failed") {
		t.Errorf("expected error log, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "vault unavailable") {
		t.Errorf("expected error message in log, got: %s", buf.String())
	}
}

func TestLoggingMiddleware_NilLogger(t *testing.T) {
	base := &stubFetcher{data: map[string]interface{}{}}
	f := middleware.Chain(base, middleware.NewLoggingMiddleware(nil))
	_, err := f.GetSecretData(context.Background(), "secret/safe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
