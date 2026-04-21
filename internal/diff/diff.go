// Package diff compares two secret snapshots and reports changes.
package diff

// ChangeType describes the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single key-level difference between two snapshots.
type Change struct {
	Key  string
	Type ChangeType
}

// Comparer compares two maps of secret data.
type Comparer struct{}

// New returns a new Comparer.
func New() *Comparer {
	return &Comparer{}
}

// Compare returns the list of changes between prev and next.
func (c *Comparer) Compare(prev, next map[string]string) []Change {
	var changes []Change

	for k, v := range next {
		if old, ok := prev[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Added})
		} else if old != v {
			changes = append(changes, Change{Key: k, Type: Changed})
		}
	}

	for k := range prev {
		if _, ok := next[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Removed})
		}
	}

	return changes
}

// HasChanges returns true if any differences exist between prev and next.
func (c *Comparer) HasChanges(prev, next map[string]string) bool {
	return len(c.Compare(prev, next)) > 0
}

// FilterByType returns only the changes that match the given ChangeType.
func (c *Comparer) FilterByType(changes []Change, t ChangeType) []Change {
	var filtered []Change
	for _, ch := range changes {
		if ch.Type == t {
			filtered = append(filtered, ch)
		}
	}
	return filtered
}
