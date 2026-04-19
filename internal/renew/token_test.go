package renew

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTokenRenewer_MissingAddress(t *testing.T) {
	_, err := NewTokenRenewer("", "token", nil)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNewTokenRenewer_MissingToken(t *testing.T) {
	_, err := NewTokenRenewer("http://vault", "", nil)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestTokenRenewer_RenewSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/token/renew-self" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Vault-Token") != "mytoken" {
			t.Error("missing vault token header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	tr, err := NewTokenRenewer(ts.URL, "mytoken", ts.Client())
	if err != nil {
		t.Fatal(err)
	}
	if err := tr.Renew(context.Background()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestTokenRenewer_RenewFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()
	tr, err := NewTokenRenewer(ts.URL, "badtoken", ts.Client())
	if err != nil {
		t.Fatal(err)
	}
	if err := tr.Renew(context.Background()); err == nil {
		t.Fatal("expected error for non-200 response")
	}
}
