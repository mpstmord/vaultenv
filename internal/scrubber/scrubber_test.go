package scrubber_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/scrubber"
)

func TestNew_DefaultPlaceholder(t *testing.T) {
	s := scrubber.New("")
	s.Add("mysecret")
	got := s.Scrub("value is mysecret here")
	want := "value is [REDACTED] here"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNew_CustomPlaceholder(t *testing.T) {
	s := scrubber.New("***")
	s.Add("tok")
	got := s.Scrub("tok is tok")
	if got != "*** is ***" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestScrub_NoSecrets(t *testing.T) {
	s := scrubber.New("")
	got := s.Scrub("nothing to hide")
	if got != "nothing to hide" {
		t.Errorf("expected unchanged text, got %q", got)
	}
}

func TestScrub_MultipleSecrets(t *testing.T) {
	s := scrubber.New("")
	s.Add("alpha", "beta")
	got := s.Scrub("alpha and beta are secrets")
	if got != "[REDACTED] and [REDACTED] are secrets" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestAdd_IgnoresEmpty(t *testing.T) {
	s := scrubber.New("")
	s.Add("", "", "real")
	if s.Len() != 1 {
		t.Errorf("expected 1 secret, got %d", s.Len())
	}
}

func TestScrubMap_RedactsValues(t *testing.T) {
	s := scrubber.New("")
	s.Add("s3cr3t")
	input := map[string]string{
		"TOKEN": "s3cr3t",
		"HOST":  "localhost",
	}
	out := s.ScrubMap(input)
	if out["TOKEN"] != "[REDACTED]" {
		t.Errorf("expected TOKEN to be redacted, got %q", out["TOKEN"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST unchanged, got %q", out["HOST"])
	}
}

func TestScrubMap_DoesNotMutateOriginal(t *testing.T) {
	s := scrubber.New("")
	s.Add("pw")
	orig := map[string]string{"PASS": "pw"}
	_ = s.ScrubMap(orig)
	if orig["PASS"] != "pw" {
		t.Error("original map was mutated")
	}
}
