package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultenv configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Secrets []SecretEntry `yaml:"secrets"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

// SecretEntry maps a Vault secret field to an environment variable.
type SecretEntry struct {
	// Mapping format: ENV_VAR=secret/path#field
	Mapping string `yaml:"mapping"`
}

// LoadFile reads and parses a YAML config file from the given path.
func LoadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}
	return &cfg, nil
}

// Validate returns an error if the configuration is incomplete.
func (c *Config) Validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("config: vault.address is required")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("config: vault.token is required")
	}
	if len(c.Secrets) == 0 {
		return fmt.Errorf("config: at least one secret mapping is required")
	}
	return nil
}
