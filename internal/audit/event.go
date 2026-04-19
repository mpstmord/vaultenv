package audit

import "time"

// EventType classifies an audit log entry.
type EventType string

const (
	EventSecretFetch  EventType = "secret_fetch"
	EventSecretInject EventType = "secret_inject"
	EventProcessExec  EventType = "process_exec"
	EventError        EventType = "error"
)

// Event represents a single auditable action performed by vaultenv.
type Event struct {
	// Timestamp is when the event occurred (UTC).
	Timestamp time.Time `json:"timestamp"`

	// Type classifies the event.
	Type EventType `json:"type"`

	// Path is the Vault secret path involved, if any.
	Path string `json:"path,omitempty"`

	// EnvVar is the environment variable name being injected, if any.
	EnvVar string `json:"env_var,omitempty"`

	// Command is the subprocess command that was executed, if any.
	Command string `json:"command,omitempty"`

	// Error holds an error message when Type == EventError.
	Error string `json:"error,omitempty"`

	// Success indicates whether the action completed without error.
	Success bool `json:"success"`
}

// NewEvent constructs an Event with the current UTC timestamp.
func NewEvent(t EventType) Event {
	return Event{
		Timestamp: time.Now().UTC(),
		Type:      t,
		Success:   true,
	}
}

// WithError marks the event as failed and records the error message.
func (e Event) WithError(err error) Event {
	e.Success = false
	if err != nil {
		e.Error = err.Error()
	}
	return e
}
