package plugin

import (
	"fmt"
	"os"
)

// EnvProvider is a built-in Provider that reads secrets from environment
// variables. The path is ignored; field is used as the variable name.
//
// This is useful for local development and CI pipelines where secrets are
// already available as environment variables.
type EnvProvider struct{}

// NewEnvProvider returns an EnvProvider ready for registration.
func NewEnvProvider() *EnvProvider { return &EnvProvider{} }

// Name returns the identifier used to select this provider.
func (e *EnvProvider) Name() string { return "env" }

// GetSecret returns the value of the environment variable named by field.
// Returns an error if the variable is not set.
func (e *EnvProvider) GetSecret(path, field string) (string, error) {
	val, ok := os.LookupEnv(field)
	if !ok {
		return "", fmt.Errorf("env provider: variable %q not set (path %q)", field, path)
	}
	return val, nil
}
