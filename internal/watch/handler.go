package watch

import (
	"fmt"
	"io"
	"strings"
)

// LogChangeHandler returns a ChangeHandler that writes a human-readable
// summary of added, removed, and updated keys to w.
func LogChangeHandler(w io.Writer) ChangeHandler {
	return func(old, new map[string]string) error {
		var lines []string

		for k := range new {
			if _, exists := old[k]; !exists {
				lines = append(lines, fmt.Sprintf("  + %s (added)", k))
			} else if old[k] != new[k] {
				lines = append(lines, fmt.Sprintf("  ~ %s (changed)", k))
			}
		}
		for k := range old {
			if _, exists := new[k]; !exists {
				lines = append(lines, fmt.Sprintf("  - %s (removed)", k))
			}
		}

		if len(lines) == 0 {
			return nil
		}
		_, err := fmt.Fprintf(w, "secret change detected:\n%s\n", strings.Join(lines, "\n"))
		return err
	}
}

// ChainChangeHandlers returns a ChangeHandler that calls each handler in
// order, stopping and returning the first non-nil error.
func ChainChangeHandlers(handlers ...ChangeHandler) ChangeHandler {
	return func(old, new map[string]string) error {
		for _, h := range handlers {
			if err := h(old, new); err != nil {
				return err
			}
		}
		return nil
	}
}
