package middleware

import (
	"context"
	"log/slog"
	"time"
)

// NewLoggingMiddleware returns a Middleware that logs each secret fetch
// with its path, duration, and any error using the provided slog.Logger.
func NewLoggingMiddleware(logger *slog.Logger) Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	return func(next Fetcher) Fetcher {
		return FetcherFunc(func(ctx context.Context, path string) (map[string]interface{}, error) {
			start := time.Now()
			data, err := next.GetSecretData(ctx, path)
			duration := time.Since(start)
			if err != nil {
				logger.Error("secret fetch failed",
					"path", path,
					"duration_ms", duration.Milliseconds(),
					"error", err.Error(),
				)
			} else {
				logger.Info("secret fetched",
					"path", path,
					"duration_ms", duration.Milliseconds(),
					"fields", len(data),
				)
			}
			return data, err
		})
	}
}
