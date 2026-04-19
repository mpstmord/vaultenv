package rotate

import (
	"log"
	"os"
	"os/exec"
)

// ReloadFunc is a function that reloads a child process.
type ReloadFunc func() error

// LogHandler returns a Handler that logs each rotation event.
func LogHandler(logger *log.Logger) Handler {
	return func(path string, _ map[string]interface{}) error {
		logger.Printf("rotate: secret updated at %s", path)
		return nil
	}
}

// ExecHandler returns a Handler that runs a shell command on rotation.
func ExecHandler(command string) Handler {
	return func(path string, _ map[string]interface{}) error {
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

// ChainHandlers combines multiple handlers; stops on first error.
func ChainHandlers(handlers ...Handler) Handler {
	return func(path string, data map[string]interface{}) error {
		for _, h := range handlers {
			if err := h(path, data); err != nil {
				return err
			}
		}
		return nil
	}
}
