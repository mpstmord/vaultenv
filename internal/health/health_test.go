package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type okChecker struct{}

func (o *okChecker) Name() string                  { return "ok" }
func (o *okChecker) Check(_ context.Context) error { return nil }

type failChecker struct{}

func (f *failChecker) Name() string                  { return "fail" }
func (f *failChecker) Check(_ context.Context) error { return errors.New("unavailable") }

func TestRunner_AllOK(t *testing.T) {
	r := NewRunner(time.Second, &okChecker{})
	rep := r.Run(context.Background())
	if rep.Status != string(StatusOK) {
		t.Fatalf("expected ok, got %s", rep.Status)
	}
	if len(rep.Checks) != 1 || rep.Checks[0].Status != StatusOK {
		t.Fatal("expected single ok check")
	}
}

func TestRunner_OneFails(t *testing.T) {
	r := NewRunner(time.Second, &okChecker{}, &failChecker{})
	rep := r.Run(context.Background())
	if rep.Status != string(StatusDegraded) {
		t.Fatalf("expected degraded, got %s", rep.Status)
	}
}

func TestRunner_DefaultTimeout(t *testing.T) {
	r := NewRunner(0, &okChecker{})
	if r.timeout != 5*time.Second {
		t.Fatalf("expected default 5s timeout, got %v", r.timeout)
	}
}

func TestVaultChecker_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	vc := NewVaultChecker(ts.URL)
	if err := vc.Check(context.Background()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestVaultChecker_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	vc := NewVaultChecker(ts.URL)
	if err := vc.Check(context.Background()); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestVaultChecker_Name(t *testing.T) {
	vc := NewVaultChecker("http://localhost:8200")
	if vc.Name() != "vault" {
		t.Fatalf("unexpected name %s", vc.Name())
	}
}
