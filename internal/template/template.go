package template

import (
	"fmt"
	"os"
	"strings"
)

// Renderer renders a template string by substituting secret references
// with resolved values from a lookup function.
type Renderer struct {
	lookup func(ref string) (string, error)
}

// NewRenderer creates a Renderer that uses the provided lookup function
// to resolve secret references of the form {{vault:path#field}}.
func NewRenderer(lookup func(ref string) (string, error)) *Renderer {
	return &Renderer{lookup: lookup}
}

// Render substitutes all {{vault:...}} placeholders in src with resolved values.
func (r *Renderer) Render(src string) (string, error) {
	var err error
	result := replacePlaceholders(src, func(ref string) string {
		if err != nil {
			return ""
		}
		val, e := r.lookup(ref)
		if e != nil {
			err = fmt.Errorf("template: resolving %q: %w", ref, e)
			return ""
		}
		return val
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// RenderFile reads a file, renders its contents, and returns the result.
func (r *Renderer) RenderFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template: reading file %q: %w", path, err)
	}
	return r.Render(string(data))
}

// replacePlaceholders replaces all {{vault:...}} tokens using fn.
func replacePlaceholders(src string, fn func(string) string) string {
	const open = "{{vault:"
	const close = "}}"
	var sb strings.Builder
	for {
		start := strings.Index(src, open)
		if start == -1 {
			sb.WriteString(src)
			break
		}
		sb.WriteString(src[:start])
		rest := src[start+len(open):]
		end := strings.Index(rest, close)
		if end == -1 {
			sb.WriteString(src[start:])
			break
		}
		ref := rest[:end]
		sb.WriteString(fn(ref))
		src = rest[end+len(close):]
	}
	return sb.String()
}
