package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	logical *vaultapi.Logical
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If Address or Token are empty, it falls back to environment variables.
func NewClient(cfg Config) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address not set: provide --vault-addr or VAULT_ADDR")
	}
	vaultCfg.Address = address

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token not set: provide --vault-token or VAULT_TOKEN")
	}
	client.SetToken(token)

	return &Client{logical: client.Logical()}, nil
}

// GetSecretData reads a KV secret at the given path and returns its data map.
func (c *Client) GetSecretData(path string) (map[string]interface{}, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	// KV v2 wraps data under secret.Data["data"]
	if data, ok := secret.Data["data"]; ok {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return secret.Data, nil
}
