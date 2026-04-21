package store

import (
	"errors"
	"sync"
)

// ErrNotFound is returned when a secret key does not exist in the store.
var ErrNotFound = errors.New("store: key not found")

// Store is a thread-safe in-memory key/value store for resolved secret values.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns an empty Store.
func New() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Set stores a secret value under the given key, overwriting any previous value.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves the value for key. Returns ErrNotFound if the key is absent.
func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Delete removes the key from the store. It is a no-op if the key does not exist.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Keys returns a snapshot of all keys currently held in the store.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries in the store.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Snapshot returns a shallow copy of all key/value pairs.
func (s *Store) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]string, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}
