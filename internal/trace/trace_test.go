package trace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/trace"
)

func TestStart_ReturnsSpanWithOperation(t *testing.T) {
	ctx, span := trace.Start(context.Background(), "fetch-secret")
	if span == nil {
		t.Fatal("expected non-nil span")
	}
	if span.Operation != "fetch-secret" {
		t.Errorf("operation = %q, want %q", span.Operation, "fetch-secret")
	}
	if span.TraceID == "" {
		t.Error("expected non-empty trace ID")
	}
	if trace.FromContext(ctx) != span {
		t.Error("span not stored in context")
	}
}

func TestStart_TraceIDIsUnique(t *testing.T) {
	_, s1 := trace.Start(context.Background(), "op")
	_, s2 := trace.Start(context.Background(), "op")
	if s1.TraceID == s2.TraceID {
		t.Error("expected unique trace IDs")
	}
}

func TestSpan_Finish_RecordsError(t *testing.T) {
	_, span := trace.Start(context.Background(), "op")
	err := errors.New("vault unavailable")
	span.Finish(err)
	if span.Err != err {
		t.Errorf("err = %v, want %v", span.Err, err)
	}
	if span.EndedAt.IsZero() {
		t.Error("EndedAt should not be zero after Finish")
	}
}

func TestSpan_Duration_Positive(t *testing.T) {
	_, span := trace.Start(context.Background(), "op")
	time.Sleep(2 * time.Millisecond)
	span.Finish(nil)
	if span.Duration() <= 0 {
		t.Errorf("expected positive duration, got %v", span.Duration())
	}
}

func TestFromContext_NoSpan_ReturnsNil(t *testing.T) {
	if s := trace.FromContext(context.Background()); s != nil {
		t.Errorf("expected nil, got %v", s)
	}
}

func TestTraceID_WithSpan(t *testing.T) {
	ctx, span := trace.Start(context.Background(), "op")
	if id := trace.TraceID(ctx); id != span.TraceID {
		t.Errorf("TraceID = %q, want %q", id, span.TraceID)
	}
}

func TestTraceID_WithoutSpan_ReturnsEmpty(t *testing.T) {
	if id := trace.TraceID(context.Background()); id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}
