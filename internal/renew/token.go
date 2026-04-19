package renew

import (
	"context"
	"errors"
	"net/http"
)

// TokenRenewer wraps an HTTP client to renew a Vault token via the API.
type TokenRenewer struct {
	address string
	token   string
	client  *http.Client
}

// NewTokenRenewer creates a TokenRenewer for the given Vault address and token.
func NewTokenRenewer(address, token string, client *http.Client) (*TokenRenewer, error) {
	if address == "" {
		return nil, errors.New("vault address is required")
	}
	if token == "" {
		return nil, errors.New("vault token is required")
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * 1e9}
	}
	return &TokenRenewer{address: address, token: token, client: client}, nil
}

// Renew calls the Vault token self-renew endpoint.
func (tr *TokenRenewer) Renew(ctx context.Context) error {
	url := tr.address + "/v1/auth/token/renew-self"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", tr.token)
	resp, err := tr.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("token renewal returned status: " + resp.Status)
	}
	return nil
}
