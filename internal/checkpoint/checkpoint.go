// Package checkpoint persists the last-seen secret version so that
// vaultenv can resume from a known state after a restart.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Record holds the persisted state for a single secret path.
type Record struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists checkpoint records to a JSON file.
type Store struct {
	mu      sync.RWMutex
	file    string
	records map[string]Record
}

// New loads an existing checkpoint file or creates an empty store.
// If the file does not exist, an empty store is returned without error.
func New(file string) (*Store, error) {
	s := &Store{
		file:    file,
		records: make(map[string]Record),
	}
	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.records); err != nil {
		return nil, err
	}
	return s, nil
}

// Set updates the record for the given path and flushes to disk.
func (s *Store) Set(path string, version int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[path] = Record{
		Path:      path,
		Version:   version,
		UpdatedAt: time.Now().UTC(),
	}
	return s.flush()
}

// Get returns the record for path and whether it was found.
func (s *Store) Get(path string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.records[path]
	return r, ok
}

// Delete removes the record for path and flushes to disk.
func (s *Store) Delete(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, path)
	return s.flush()
}

// flush writes the current records to disk. Caller must hold the write lock.
func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0o600)
}
