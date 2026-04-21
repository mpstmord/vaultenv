package store

import (
	"context"
	"fmt"
)

// Fetcher retrieves secret data for a given path.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Mapping describes how a single Vault field maps to an environment variable.
type Mapping struct {
	// EnvKey is the environment variable name to populate.
	EnvKey string
	// Path is the Vault secret path.
	Path string
	// Field is the key within the secret's data map.
	Field string
}

// Loader fetches secrets from Vault and populates a Store.
type Loader struct {
	fetcher  Fetcher
	store    *Store
	mappings []Mapping
}

// NewLoader creates a Loader that will resolve mappings using fetcher and
// write results into store.
func NewLoader(fetcher Fetcher, store *Store, mappings []Mapping) *Loader {
	return &Loader{
		fetcher:  fetcher,
		store:    store,
		mappings: mappings,
	}
}

// Load resolves all mappings and stores the resulting values. It returns the
// first error encountered, leaving previously resolved values in the store.
func (l *Loader) Load(ctx context.Context) error {
	for _, m := range l.mappings {
		data, err := l.fetcher.GetSecretData(ctx, m.Path)
		if err != nil {
			return fmt.Errorf("store: fetch %q: %w", m.Path, err)
		}
		raw, ok := data[m.Field]
		if !ok {
			return fmt.Errorf("store: field %q not found in secret %q", m.Field, m.Path)
		}
		val, ok := raw.(string)
		if !ok {
			return fmt.Errorf("store: field %q in %q is not a string", m.Field, m.Path)
		}
		l.store.Set(m.EnvKey, val)
	}
	return nil
}
