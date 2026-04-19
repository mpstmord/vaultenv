package env

import (
	"fmt"

	"github.com/your-org/vaultenv/internal/vault"
)

// Resolver fetches secret values from Vault for a list of mappings.
type Resolver struct {
	client *vault.Client
}

// NewResolver creates a new Resolver backed by the given Vault client.
func NewResolver(client *vault.Client) *Resolver {
	return &Resolver{client: client}
}

// Resolve fetches all secret values for the given mappings and returns a map
// of environment variable name to secret value.
func (r *Resolver) Resolve(mappings []*SecretMapping) (map[string]string, error) {
	result := make(map[string]string, len(mappings))
	for _, m := range mappings {
		data, err := r.client.GetSecretData(m.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secret at %q: %w", m.Path, err)
		}
		val, ok := data[m.Field]
		if !ok {
			return nil, fmt.Errorf("field %q not found in secret at %q", m.Field, m.Path)
		}
		strVal, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("field %q in secret at %q is not a string", m.Field, m.Path)
		}
		result[m.EnvVar] = strVal
	}
	return result, nil
}
