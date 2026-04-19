// Package audit provides structured audit logging for secret access events.
package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventSecretFetch  EventType = "secret_fetch"
	EventSecretResolve EventType = "secret_resolve"
	EventExec         EventType = "exec"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time  `json:"timestamp"`
	Type      EventType  `json:"type"`
	Path      string     `json:"path,omitempty"`
	Field     string     `json:"field,omitempty"`
	Command   string     `json:"command,omitempty"`
	Success   bool       `json:"success"`
	Error     string     `json:"error,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w       io.Writer
	enabled bool
}

// NewLogger creates a Logger writing to w. If w is nil, logging is a no-op.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		return &Logger{enabled: false}
	}
	return &Logger{w: w, enabled: true}
}

// NewFileLogger opens path for appending and returns a Logger backed by it.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	return NewLogger(f), nil
}

// Log emits an audit event. It is safe to call on a nil or disabled Logger.
func (l *Logger) Log(e Event) error {
	if l == nil || !l.enabled {
		return nil
	}
	e.Timestamp = time.Now().UTC()
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.w.Write(data)
	return err
}
