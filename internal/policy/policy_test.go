package policy

import (
	"testing"
)

func TestNew_EmptyPathReturnsError(t *testing.T) {
	_, err := New([]Rule{{Path: "", Allow: true}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_Valid(t *testing.T) {
	p, err := New([]Rule{{Path: "secret/*", Allow: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 rule, got %d", p.Len())
	}
}

func TestAllow_NoRules_Denied(t *testing.T) {
	p, _ := New(nil)
	if p.Allow("secret/app/db") {
		t.Fatal("expected denial with no rules")
	}
}

func TestAllow_ExactMatch(t *testing.T) {
	p, _ := New([]Rule{{Path: "secret/app/db", Allow: true}})
	if !p.Allow("secret/app/db") {
		t.Fatal("expected allow for exact match")
	}
	if p.Allow("secret/app/other") {
		t.Fatal("expected denial for non-matching path")
	}
}

func TestAllow_WildcardMatch(t *testing.T) {
	p, _ := New([]Rule{{Path: "secret/app/*", Allow: true}})
	if !p.Allow("secret/app/db") {
		t.Fatal("expected allow for wildcard match")
	}
	if p.Allow("secret/app/nested/deep") {
		t.Fatal("expected denial for nested path beyond wildcard")
	}
}

func TestAllow_DoubleStarMatchesAll(t *testing.T) {
	p, _ := New([]Rule{{Path: "**", Allow: true}})
	if !p.Allow("anything/goes/here") {
		t.Fatal("expected allow for ** pattern")
	}
}

func TestAllow_FirstMatchWins(t *testing.T) {
	p, _ := New([]Rule{
		{Path: "secret/app/admin", Allow: false},
		{Path: "secret/app/*", Allow: true},
	})
	if p.Allow("secret/app/admin") {
		t.Fatal("expected denial: first rule blocks admin")
	}
	if !p.Allow("secret/app/config") {
		t.Fatal("expected allow: second rule permits other paths")
	}
}

func TestAllow_DenyRule(t *testing.T) {
	p, _ := New([]Rule{{Path: "secret/*", Allow: false}})
	if p.Allow("secret/db") {
		t.Fatal("expected denial from explicit deny rule")
	}
}
