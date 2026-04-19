// Package process handles spawning child processes with injected environments.
package process

import (
	"errors"
	"os"
	"os/exec"
)

// Runner executes a command with a given environment.
type Runner struct {
	Env []string
}

// NewRunner creates a Runner with the provided environment slice.
func NewRunner(env []string) *Runner {
	return &Runner{Env: env}
}

// Run executes the given command and arguments, inheriting stdin/stdout/stderr.
// The process environment is replaced with Runner.Env.
func (r *Runner) Run(command string, args []string) error {
	if command == "" {
		return errors.New("process: command must not be empty")
	}

	path, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	cmd := exec.Command(path, args...)
	cmd.Env = r.Env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// MergeEnv merges base environment variables with overrides.
// Keys present in overrides replace those in base.
func MergeEnv(base []string, overrides map[string]string) []string {
	result := make([]string, 0, len(base)+len(overrides))
	seen := make(map[string]bool, len(overrides))

	for key := range overrides {
		seen[key] = false
	}

	for _, entry := range base {
		key := envKey(entry)
		if _, isOverride := overrides[key]; isOverride {
			seen[key] = true
			result = append(result, key+"="+overrides[key])
		} else {
			result = append(result, entry)
		}
	}

	for key, appended := range seen {
		if !appended {
			result = append(result, key+"="+overrides[key])
		}
	}

	return result
}

func envKey(entry string) string {
	for i, c := range entry {
		if c == '=' {
			return entry[:i]
		}
	}
	return entry
}
