package renew

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockLogger struct {
	mu   sync.Mutex
	msgs []string
	errs []string
}

func (m *mockLogger) Info(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, msg)
}

func (m *mockLogger) Error(msg string, _ error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errs = append(m.errs, msg)
}

func TestRenewer_CallsRenewFunc(t *testing.T) {
	var count int
	var mu sync.Mutex
	log := &mockLogger{}
	r := NewRenewer(20*time.Millisecond, func(_ context.Context) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	}, log)
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	r.Start(ctx)
	mu.Lock()
	defer mu.Unlock()
	if count < 2 {
		t.Errorf("expected at least 2 renewals, got %d", count)
	}
}

func TestRenewer_LogsError(t *testing.T) {
	log := &mockLogger{}
	r := NewRenewer(20*time.Millisecond, func(_ context.Context) error {
		return errors.New("vault unavailable")
	}, log)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	r.Start(ctx)
	log.mu.Lock()
	defer log.mu.Unlock()
	if len(log.errs) == 0 {
		t.Error("expected error to be logged")
	}
}

func TestNewRenewer_DefaultInterval(t *testing.T) {
	r := NewRenewer(0, func(_ context.Context) error { return nil }, &mockLogger{})
	if r.interval != 10*time.Minute {
		t.Errorf("expected default interval 10m, got %v", r.interval)
	}
}
