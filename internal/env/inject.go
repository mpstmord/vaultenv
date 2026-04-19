package env

import (
	"fmt"
	"os"
	"strings"
)

// SecretMapping maps an environment variable name to a Vault path and field.
type SecretMapping struct {
	EnvVar string
	Path   string
	Field  string
}

// ParseMapping parses a mapping string of the form ENV_VAR=vault/path#field.
func ParseMapping(s string) (*SecretMapping, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid mapping %q: expected ENV_VAR=path#field", s)
	}
	envVar := strings.TrimSpace(parts[0])
	if envVar == "" {
		return nil, fmt.Errorf("invalid mapping %q: empty env var name", s)
	}
	right := strings.SplitN(parts[1], "#", 2)
	if len(right) != 2 {
		return nil, fmt.Errorf("invalid mapping %q: expected path#field", s)
	}
	path := strings.TrimSpace(right[0])
	field := strings.TrimSpace(right[1])
	if path == "" || field == "" {
		return nil, fmt.Errorf("invalid mapping %q: path and field must not be empty", s)
	}
	return &SecretMapping{EnvVar: envVar, Path: path, Field: field}, nil
}

// InjectIntoEnv sets environment variables from the provided resolved values.
func InjectIntoEnv(resolved map[string]string) error {
	for k, v := range resolved {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set env var %q: %w", k, err)
		}
	}
	return nil
}
