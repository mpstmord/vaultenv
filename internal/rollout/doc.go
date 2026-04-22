// Package rollout implements a gradual secret rollout gate for vaultenv.
//
// It allows operators to incrementally expose a new secret version to a
// controlled percentage of processes, using a deterministic hash strategy
// so that the same key always resolves consistently within a given run.
package rollout
