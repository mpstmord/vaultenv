// Package plugin provides a registry for secret backend plugins,
// allowing vaultenv to support alternative secret sources beyond HashiCorp Vault.
package plugin

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNotFound is returned when a plugin is not registered under the given name.
var ErrNotFound = errors.New("plugin not found")

// ErrAlreadyRegistered is returned when a plugin name is registered more than once.
var ErrAlreadyRegistered = errors.New("plugin already registered")

// Provider is the interface that all secret backend plugins must implement.
type Provider interface {
	// Name returns the unique identifier for this provider.
	Name() string

	// GetSecret retrieves a secret value by path and field.
	GetSecret(path, field string) (string, error)
}

// Registry holds registered secret backend plugins.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Provider
}

// New creates and returns an empty Registry.
func New() *Registry {
	return &Registry{
		plugins: make(map[string]Provider),
	}
}

// Register adds a Provider to the registry.
// Returns ErrAlreadyRegistered if a provider with the same name exists.
func (r *Registry) Register(p Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := p.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("%w: %s", ErrAlreadyRegistered, name)
	}
	r.plugins[name] = p
	return nil
}

// Get retrieves a registered Provider by name.
// Returns ErrNotFound if no provider is registered under that name.
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return p, nil
}

// Names returns a sorted list of all registered provider names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

// Unregister removes a provider by name. Returns ErrNotFound if absent.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.plugins[name]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	delete(r.plugins, name)
	return nil
}
