package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestGetSecretData_KVv2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"API_KEY":"supersecret"}}}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := client.GetSecretData("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}

	if data["API_KEY"] != "supersecret" {
		t.Errorf("expected API_KEY=supersecret, got %v", data["API_KEY"])
	}
}

func TestGetSecretData_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.GetSecretData("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
