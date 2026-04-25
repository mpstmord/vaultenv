// Package lock provides per-key advisory locking for vaultenv operations.
//
// It prevents concurrent goroutines from simultaneously fetching or writing
// the same secret key, reducing redundant Vault API calls and avoiding
// race conditions during cache population.
//
// Usage:
//
//	l := lock.New(2 * time.Second)
//	if err := l.Lock(ctx, path); err != nil {
//		return err
//	}
//	defer l.Unlock(path)
package lock
