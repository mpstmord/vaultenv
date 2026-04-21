// Package trace provides request tracing for Vault secret fetches.
// Each fetch operation is assigned a unique trace ID that propagates
// through the middleware stack for correlation in audit logs.
package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type contextKey struct{}

// Span holds tracing metadata for a single secret fetch operation.
type Span struct {
	TraceID   string
	Operation string
	StartedAt time.Time
	EndedAt   time.Time
	Err       error
}

// Duration returns the elapsed time of the span.
func (s *Span) Duration() time.Duration {
	return s.EndedAt.Sub(s.StartedAt)
}

// Finish marks the span as complete, recording the end time and any error.
func (s *Span) Finish(err error) {
	s.EndedAt = time.Now()
	s.Err = err
}

// Start creates a new Span and returns a context carrying it.
func Start(ctx context.Context, operation string) (context.Context, *Span) {
	span := &Span{
		TraceID:   newTraceID(),
		Operation: operation,
		StartedAt: time.Now(),
	}
	return context.WithValue(ctx, contextKey{}, span), span
}

// FromContext retrieves the Span stored in ctx, or nil if none.
func FromContext(ctx context.Context) *Span {
	v := ctx.Value(contextKey{})
	if v == nil {
		return nil
	}
	s, _ := v.(*Span)
	return s
}

// TraceID returns the trace ID from ctx, or an empty string.
func TraceID(ctx context.Context) string {
	if s := FromContext(ctx); s != nil {
		return s.TraceID
	}
	return ""
}

func newTraceID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
