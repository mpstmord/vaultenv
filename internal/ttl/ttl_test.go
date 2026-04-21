package ttl_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/ttl"
)

func TestNewLease_ValidBeforeExpiry(t *testing.T) {
	l := ttl.NewLease("s3cr3t", 5*time.Second)
	val, err := l.Value()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "s3cr3t" {
		t.Fatalf("expected 's3cr3t', got %q", val)
	}
}

func TestNewLease_ExpiredImmediately(t *testing.T) {
	l := ttl.NewLease("s3cr3t", -1*time.Second)
	_, err := l.Value()
	if err != ttl.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestLease_TTL_Positive(t *testing.T) {
	l := ttl.NewLease("x", 10*time.Second)
	if l.TTL() <= 0 {
		t.Fatal("expected positive TTL")
	}
}

func TestLease_TTL_Negative(t *testing.T) {
	l := ttl.NewLease("x", -1*time.Second)
	if l.TTL() > 0 {
		t.Fatal("expected non-positive TTL for expired lease")
	}
}

func TestLease_Expired_False(t *testing.T) {
	l := ttl.NewLease("x", 5*time.Second)
	if l.Expired() {
		t.Fatal("lease should not be expired yet")
	}
}

func TestLease_Expired_True(t *testing.T) {
	l := ttl.NewLease("x", -1*time.Millisecond)
	if !l.Expired() {
		t.Fatal("lease should be expired")
	}
}

func TestLease_Renew_ExtendsValidity(t *testing.T) {
	l := ttl.NewLease("x", -1*time.Second)
	if !l.Expired() {
		t.Fatal("lease should start expired")
	}
	l.Renew(5 * time.Second)
	if l.Expired() {
		t.Fatal("lease should be valid after renew")
	}
	val, err := l.Value()
	if err != nil {
		t.Fatalf("expected no error after renew, got %v", err)
	}
	if val != "x" {
		t.Fatalf("expected 'x', got %q", val)
	}
}
