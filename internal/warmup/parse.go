package warmup

import (
	"fmt"
	"strings"
)

// ParseMappings parses a slice of strings in the form
// "ENV_KEY=secret/path#field" into Mapping values.
func ParseMappings(specs []string) ([]Mapping, error) {
	out := make([]Mapping, 0, len(specs))
	for _, spec := range specs {
		m, err := parseOne(spec)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func parseOne(spec string) (Mapping, error) {
	eqIdx := strings.IndexByte(spec, '=')
	if eqIdx < 1 {
		return Mapping{}, fmt.Errorf("warmup: invalid mapping %q: missing '='" , spec)
	}
	envKey := spec[:eqIdx]
	rest := spec[eqIdx+1:]

	hashIdx := strings.LastIndexByte(rest, '#')
	if hashIdx < 1 {
		return Mapping{}, fmt.Errorf("warmup: invalid mapping %q: missing '#'", spec)
	}
	path := rest[:hashIdx]
	field := rest[hashIdx+1:]

	if path == "" {
		return Mapping{}, fmt.Errorf("warmup: invalid mapping %q: empty path", spec)
	}
	if field == "" {
		return Mapping{}, fmt.Errorf("warmup: invalid mapping %q: empty field", spec)
	}

	return Mapping{Path: path, Field: field, EnvKey: envKey}, nil
}
