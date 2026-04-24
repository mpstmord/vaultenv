// Package notify provides webhook and channel-based notification delivery
// for secret rotation and health events in vaultenv.
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload is the JSON body sent to a webhook endpoint.
type Payload struct {
	Event     string            `json:"event"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}

// Notifier sends notifications to one or more configured destinations.
type Notifier struct {
	client   *http.Client
	endpoints []string
}

// Option configures a Notifier.
type Option func(*Notifier)

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(n *Notifier) {
		n.client = c
	}
}

// New creates a Notifier that will deliver payloads to the given webhook URLs.
// Returns an error if no endpoints are provided.
func New(endpoints []string, opts ...Option) (*Notifier, error) {
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("notify: at least one endpoint is required")
	}
	n := &Notifier{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoints: endpoints,
	}
	for _, o := range opts {
		o(n)
	}
	return n, nil
}

// Send delivers a notification payload to every configured endpoint.
// It collects all errors and returns a combined error if any delivery fails.
func (n *Notifier) Send(ctx context.Context, event string, details map[string]string) error {
	p := Payload{
		Event:     event,
		Timestamp: time.Now().UTC(),
		Details:   details,
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	var errs []error
	for _, ep := range n.endpoints {
		if err := n.post(ctx, ep, body); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", ep, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify: delivery failures: %v", errs)
	}
	return nil
}

// post sends a single HTTP POST request to the given URL.
func (n *Notifier) post(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}
