package filter_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/filter"
)

func TestAllow_NoPatterns(t *testing.T) {
	f := filter.New(nil, nil)
	if !f.Allow("ANY_KEY") {
		t.Error("expected all keys allowed when no patterns set")
	}
}

func TestAllow_IncludePattern(t *testing.T) {
	f := filter.New([]string{"DB_*"}, nil)
	if !f.Allow("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be included")
	}
	if f.Allow("AWS_SECRET") {
		t.Error("expected AWS_SECRET to be excluded")
	}
}

func TestAllow_ExcludePattern(t *testing.T) {
	f := filter.New(nil, []string{"AWS_*"})
	if f.Allow("AWS_SECRET_KEY") {
		t.Error("expected AWS_SECRET_KEY to be excluded")
	}
	if !f.Allow("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]string{"DB_*"}, []string{"DB_PASSWORD"})
	if f.Allow("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD excluded despite include pattern")
	}
	if !f.Allow("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f := filter.New([]string{"db_*"}, nil)
	if !f.Allow("DB_HOST") {
		t.Error("expected case-insensitive match")
	}
}

func TestApply_FiltersMap(t *testing.T) {
	f := filter.New([]string{"APP_*"}, nil)
	input := map[string]string{
		"APP_PORT": "8080",
		"DB_HOST":  "localhost",
		"APP_ENV":  "prod",
	}
	out := f.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be filtered out")
	}
}
