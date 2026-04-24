package batch

import (
	"context"
	"log"
)

// LoggingFetcher wraps a Fetcher and logs each request and any errors.
type LoggingFetcher struct {
	upstream Fetcher
	logger   *log.Logger
}

// NewLoggingFetcher returns a Fetcher that logs every GetSecretData call.
func NewLoggingFetcher(upstream Fetcher, logger *log.Logger) *LoggingFetcher {
	if upstream == nil {
		panic("batch: upstream fetcher must not be nil")
	}
	if logger == nil {
		logger = log.Default()
	}
	return &LoggingFetcher{upstream: upstream, logger: logger}
}

// GetSecretData delegates to the upstream fetcher and logs the outcome.
func (l *LoggingFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	l.logger.Printf("batch: fetching path=%s", path)
	data, err := l.upstream.GetSecretData(ctx, path)
	if err != nil {
		l.logger.Printf("batch: error path=%s err=%v", path, err)
	}
	return data, err
}
