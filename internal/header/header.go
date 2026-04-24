// Package header provides utilities for injecting custom HTTP headers
// into Vault client requests, such as X-Request-ID or X-Vault-Namespace.
package header

import (
	"fmt"
	"net/http"
	"strings"
)

// Provider holds a set of static headers to be applied to HTTP requests.
type Provider struct {
	headers map[string]string
}

// New creates a Provider from a slice of "Key: Value" strings.
// Returns an error if any entry is malformed.
func New(pairs []string) (*Provider, error) {
	h := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.Index(p, ":")
		if idx <= 0 {
			return nil, fmt.Errorf("header: malformed pair %q: expected \"Key: Value\"", p)
		}
		key := strings.TrimSpace(p[:idx])
		val := strings.TrimSpace(p[idx+1:])
		if key == "" {
			return nil, fmt.Errorf("header: empty key in pair %q", p)
		}
		h[key] = val
	}
	return &Provider{headers: h}, nil
}

// Apply sets all stored headers on the given HTTP request.
// Existing values for the same key are replaced.
func (p *Provider) Apply(req *http.Request) {
	if req == nil {
		return
	}
	for k, v := range p.headers {
		req.Header.Set(k, v)
	}
}

// Headers returns a copy of the stored headers map.
func (p *Provider) Headers() map[string]string {
	out := make(map[string]string, len(p.headers))
	for k, v := range p.headers {
		out[k] = v
	}
	return out
}

// Len returns the number of stored headers.
func (p *Provider) Len() int {
	return len(p.headers)
}
