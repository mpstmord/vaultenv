package shadow_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/your-org/vaultenv/internal/shadow"
)

// staticFetcher returns a fixed map or error.
type staticFetcher struct {
	data map[string]interface{}
	err  error
}

func (f *staticFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return f.data, f.err
}

// captureLogger records log lines.
type captureLogger struct {
	mu   sync.Mutex
	lines []string
}

func (l *captureLogger) Printf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *captureLogger) contains(sub string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, line := range l.lines {
		if strings.Contains(line, sub) {
			return true
		}
	}
	return false
}

func TestNew_PanicsOnNilPrimary(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil primary")
		}
	}()
	shadow.New(nil, &staticFetcher{}, nil)
}

func TestNew_PanicsOnNilSecondary(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil secondary")
		}
	}()
	shadow.New(&staticFetcher{}, nil, nil)
}

func TestGetSecretData_ReturnsPrimaryResult(t *testing.T) {
	primary := &staticFetcher{data: map[string]interface{}{"key": "primary"}}
	secondary := &staticFetcher{data: map[string]interface{}{"key": "secondary"}}
	log := &captureLogger{}

	s := shadow.New(primary, secondary, log)
	data, err := s.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "primary" {
		t.Errorf("expected primary value, got %v", data["key"])
	}
}

func TestGetSecretData_LogsMismatch(t *testing.T) {
	primary := &staticFetcher{data: map[string]interface{}{"key": "v1"}}
	secondary := &staticFetcher{data: map[string]interface{}{"key": "v2"}}
	log := &captureLogger{}

	s := shadow.New(primary, secondary, log)
	_, _ = s.GetSecretData(context.Background(), "secret/test")

	if !log.contains("mismatch") {
		t.Error("expected mismatch log entry")
	}
}

func TestGetSecretData_LogsSecondaryError(t *testing.T) {
	primary := &staticFetcher{data: map[string]interface{}{"key": "v1"}}
	secondary := &staticFetcher{err: errors.New("secondary down")}
	log := &captureLogger{}

	s := shadow.New(primary, secondary, log)
	data, err := s.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected primary error: %v", err)
	}
	if data["key"] != "v1" {
		t.Errorf("expected primary data")
	}
	if !log.contains("secondary error") {
		t.Error("expected secondary error log")
	}
}

func TestGetSecretData_NoLogOnMatch(t *testing.T) {
	shared := map[string]interface{}{"key": "same"}
	primary := &staticFetcher{data: shared}
	secondary := &staticFetcher{data: map[string]interface{}{"key": "same"}}
	log := &captureLogger{}

	s := shadow.New(primary, secondary, log)
	_, _ = s.GetSecretData(context.Background(), "secret/test")

	if len(log.lines) > 0 {
		t.Errorf("expected no log entries, got: %v", log.lines)
	}
}

func TestGetSecretData_PropagatesPrimaryError(t *testing.T) {
	primary := &staticFetcher{err: errors.New("vault unavailable")}
	secondary := &staticFetcher{data: map[string]interface{}{"key": "v1"}}
	log := &captureLogger{}

	s := shadow.New(primary, secondary, log)
	_, err := s.GetSecretData(context.Background(), "secret/test")
	if err == nil {
		t.Fatal("expected primary error to propagate")
	}
}
