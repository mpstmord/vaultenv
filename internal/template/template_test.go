package template

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func fixedLookup(vals map[string]string) func(string) (string, error) {
	return func(ref string) (string, error) {
		v, ok := vals[ref]
		if !ok {
			return "", errors.New("not found: " + ref)
		}
		return v, nil
	}
}

func TestRender_NoPlaceholders(t *testing.T) {
	r := NewRenderer(fixedLookup(nil))
	out, err := r.Render("hello world")
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello world" {
		t.Errorf("expected 'hello world', got %q", out)
	}
}

func TestRender_SinglePlaceholder(t *testing.T) {
	r := NewRenderer(fixedLookup(map[string]string{"secret/app#password": "s3cr3t"}))
	out, err := r.Render("pass={{vault:secret/app#password}}")
	if err != nil {
		t.Fatal(err)
	}
	if out != "pass=s3cr3t" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_MultiplePlaceholders(t *testing.T) {
	vals := map[string]string{
		"secret/db#user": "admin",
		"secret/db#pass": "hunter2",
	}
	r := NewRenderer(fixedLookup(vals))
	out, err := r.Render("{{vault:secret/db#user}}:{{vault:secret/db#pass}}")
	if err != nil {
		t.Fatal(err)
	}
	if out != "admin:hunter2" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_LookupError(t *testing.T) {
	r := NewRenderer(fixedLookup(nil))
	_, err := r.Render("val={{vault:missing/path#field}}")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRender_EmptyString(t *testing.T) {
	r := NewRenderer(fixedLookup(nil))
	out, err := r.Render("")
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestRenderFile_Success(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	_ = os.WriteFile(p, []byte("token={{vault:sec#tok}}"), 0600)

	r := NewRenderer(fixedLookup(map[string]string{"sec#tok": "abc123"}))
	out, err := r.RenderFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if out != "token=abc123" {
		t.Errorf("unexpected: %q", out)
	}
}

func TestRenderFile_Missing(t *testing.T) {
	r := NewRenderer(fixedLookup(nil))
	_, err := r.RenderFile("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
