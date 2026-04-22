// Package rollout provides a gradual secret rollout mechanism that controls
// what percentage of processes receive a new secret version.
package rollout

import (
	"errors"
	"hash/fnv"
	"sync"
)

// Strategy determines how rollout percentage is evaluated.
type Strategy int

const (
	// StrategyHash uses a deterministic hash of the key to decide inclusion.
	StrategyHash Strategy = iota
	// StrategyRandom uses a random value per call (non-deterministic).
	StrategyRandom
)

// Config holds the rollout configuration.
type Config struct {
	// Percentage is the fraction of traffic (0–100) that receives the new value.
	Percentage int
	Strategy   Strategy
}

// Gate decides whether a given key should receive the new secret value.
type Gate struct {
	mu  sync.RWMutex
	cfg Config
}

// New creates a Gate with the provided config.
// Percentage must be in [0, 100].
func New(cfg Config) (*Gate, error) {
	if cfg.Percentage < 0 || cfg.Percentage > 100 {
		return nil, errors.New("rollout: percentage must be between 0 and 100")
	}
	return &Gate{cfg: cfg}, nil
}

// SetPercentage updates the rollout percentage at runtime.
func (g *Gate) SetPercentage(pct int) error {
	if pct < 0 || pct > 100 {
		return errors.New("rollout: percentage must be between 0 and 100")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cfg.Percentage = pct
	return nil
}

// Allow reports whether the given key should receive the new secret value.
func (g *Gate) Allow(key string) bool {
	g.mu.RLock()
	pct := g.cfg.Percentage
	strategy := g.cfg.Strategy
	g.mu.RUnlock()

	if pct == 0 {
		return false
	}
	if pct == 100 {
		return true
	}

	switch strategy {
	case StrategyHash:
		h := fnv.New32a()
		_, _ = h.Write([]byte(key))
		return int(h.Sum32()%100) < pct
	default:
		// StrategyRandom: not deterministic; use hash as default fallback.
		h := fnv.New32a()
		_, _ = h.Write([]byte(key))
		return int(h.Sum32()%100) < pct
	}
}
