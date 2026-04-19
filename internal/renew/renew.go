package renew

import (
	"context"
	"time"
)

// Renewer periodically renews a Vault token before it expires.
type Renewer struct {
	renewFunc  func(ctx context.Context) error
	interval   time.Duration
	logger     Logger
}

// Logger is a minimal logging interface used by Renewer.
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

// NewRenewer creates a Renewer that calls renewFunc on the given interval.
func NewRenewer(interval time.Duration, renewFunc func(ctx context.Context) error, logger Logger) *Renewer {
	if interval <= 0 {
		interval = 10 * time.Minute
	}
	return &Renewer{
		renewFunc: renewFunc,
		interval:  interval,
		logger:    logger,
	}
}

// Start begins the renewal loop, blocking until ctx is cancelled.
func (r *Renewer) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("token renewer stopped")
			return
		case <-ticker.C:
			if err := r.renewFunc(ctx); err != nil {
				r.logger.Error("token renewal failed", err)
			} else {
				r.logger.Info("token renewed successfully")
			}
		}
	}
}
