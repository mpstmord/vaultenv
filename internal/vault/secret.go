package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretPath represents a parsed Vault secret path with optional mount and key.
type SecretPath struct {
	Mount string
	Path  string
}

// ParseSecretPath parses a secret path in the format "mount/path/to/secret".
// The first segment is treated as the mount point.
func ParseSecretPath(raw string) (SecretPath, error) {
	raw = strings.TrimPrefix(raw, "/")
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return SecretPath{}, fmt.Errorf("invalid secret path %q: expected <mount>/<path>", raw)
	}
	return SecretPath{Mount: parts[0], Path: parts[1]}, nil
}

// GetSecretData retrieves the data map for a KV v2 secret at the given path.
// path should be in the format "<mount>/<secret-path>".
func (c *Client) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	sp, err := ParseSecretPath(path)
	if err != nil {
		return nil, err
	}

	secret, err := c.vault.KVv2(sp.Mount).Get(ctx, sp.Path)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to read secret %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: secret %q not found or empty", path)
	}
	return secret.Data, nil
}
