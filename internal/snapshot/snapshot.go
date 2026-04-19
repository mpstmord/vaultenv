// Package snapshot captures and compares secret environment snapshots.
package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

// Snapshot holds a point-in-time copy of resolved secret key/value pairs.
type Snapshot struct {
	data map[string]string
}

// New creates a Snapshot from the provided map, copying the data.
func New(data map[string]string) *Snapshot {
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return &Snapshot{data: copy}
}

// Get returns the value for key and whether it was present.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

// Keys returns a sorted list of all keys.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Digest returns a stable SHA-256 hex digest of the snapshot contents.
func (s *Snapshot) Digest() string {
	type kv struct {
		K, V string
	}
	pairs := make([]kv, 0, len(s.data))
	for k, v := range s.data {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].K < pairs[j].K })
	b, _ := json.Marshal(pairs)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Equal reports whether two snapshots have identical contents.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if s == nil || other == nil {
		return s == other
	}
	return s.Digest() == other.Digest()
}

// Diff returns keys whose values differ between s and other,
// including keys added or removed.
func (s *Snapshot) Diff(other *Snapshot) []string {
	seen := map[string]struct{}{}
	var changed []string

	for k, v := range s.data {
		seen[k] = struct{}{}
		if ov, ok := other.data[k]; !ok || ov != v {
			changed = append(changed, k)
		}
	}
	for k := range other.data {
		if _, ok := seen[k]; !ok {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}
