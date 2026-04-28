package observe

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"
)

type logEntry struct {
	Time     string `json:"time"`
	Path     string `json:"path"`
	DurationMs int64  `json:"duration_ms"`
	Cached   bool   `json:"cached"`
	Error    string `json:"error,omitempty"`
}

// LogObserver returns an Observer that writes a JSON line to w for every
// fetch event. If w is nil, os.Stderr is used.
func LogObserver(w io.Writer) Observer {
	if w == nil {
		w = os.Stderr
	}
	enc := json.NewEncoder(w)
	return func(_ context.Context, e Event) {
		entry := logEntry{
			Time:       time.Now().UTC().Format(time.RFC3339),
			Path:       e.Path,
			DurationMs: e.Duration.Milliseconds(),
			Cached:     e.Cached,
		}
		if e.Err != nil {
			entry.Error = e.Err.Error()
		}
		// Ignore encode errors — observers must not disrupt the caller.
		_ = enc.Encode(entry)
	}
}
