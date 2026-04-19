package redact

import (
	"testing"
)

func TestScrub_NoSecrets(t *testing.T) {
	r := New(nil)
	got := r.Scrub("hello world")
	if got != "hello world" {
		t.Fatalf("expected unchanged text, got %q", got)
	}
}

func TestScrub_SingleSecret(t *testing.T) {
	r := New([]string{"s3cr3t"})
	got := r.Scrub("password is s3cr3t ok")
	want := "password is [REDACTED] ok"
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}

func TestScrub_MultipleOccurrences(t *testing.T) {
	r := New([]string{"tok"})
	got := r.Scrub("tok and tok again")
	want := "[REDACTED] and [REDACTED] again"
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}

func TestScrub_EmptySecretIgnored(t *testing.T) {
	r := New([]string{"", "real"})
	if len(r.secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(r.secrets))
	}
}

func TestScrubMap_RedactsValues(t *testing.T) {
	r := New([]string{"hunter2"})
	env := map[string]string{
		"PASSWORD": "hunter2",
		"USER":     "alice",
	}
	out := r.ScrubMap(env)
	if out["PASSWORD"] != "[REDACTED]" {
		t.Fatalf("expected redacted password, got %q", out["PASSWORD"])
	}
	if out["USER"] != "alice" {
		t.Fatalf("expected unchanged user, got %q", out["USER"])
	}
}

func TestScrubMap_KeysUnchanged(t *testing.T) {
	r := New([]string{"PASSWORD"})
	env := map[string]string{"PASSWORD": "PASSWORD"}
	out := r.ScrubMap(env)
	if _, ok := out["PASSWORD"]; !ok {
		t.Fatal("key should not be redacted")
	}
}

func TestAdd_AppendsSecrets(t *testing.T) {
	r := New([]string{"first"})
	r.Add("second")
	got := r.Scrub("first second")
	want := "[REDACTED] [REDACTED]"
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}
