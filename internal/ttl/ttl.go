// Package ttl provides a time-to-live policy for secret leases.
package ttl

import (
	"errors"
	"time"
)

// ErrExpired is returned when a lease has passed its deadline.
var ErrExpired = errors.New("ttl: lease expired")

// Lease tracks the expiry of a single secret value.
type Lease struct {
	value     string
	expiresAt time.Time
}

// NewLease creates a Lease that expires after the given duration.
// A zero or negative duration causes the lease to expire immediately.
func NewLease(value string, d time.Duration) *Lease {
	return &Lease{
		value:     value,
		expiresAt: time.Now().Add(d),
	}
}

// Value returns the secret value if the lease is still valid, or ErrExpired.
func (l *Lease) Value() (string, error) {
	if time.Now().After(l.expiresAt) {
		return "", ErrExpired
	}
	return l.value, nil
}

// TTL returns the remaining duration of the lease.
// A non-positive value means the lease has expired.
func (l *Lease) TTL() time.Duration {
	return time.Until(l.expiresAt)
}

// Expired reports whether the lease has expired.
func (l *Lease) Expired() bool {
	return time.Now().After(l.expiresAt)
}

// Renew extends the lease by the given duration from now.
func (l *Lease) Renew(d time.Duration) {
	l.expiresAt = time.Now().Add(d)
}
