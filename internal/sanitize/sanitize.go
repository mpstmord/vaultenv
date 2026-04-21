// Package sanitize provides utilities for cleaning and validating
// environment variable names and secret values before injection.
package sanitize

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidEnvKey is returned when an env key contains illegal characters.
	ErrInvalidEnvKey = errors.New("invalid environment variable name")
	// ErrEmptyValue is returned when a secret value is empty.
	ErrEmptyValue = errors.New("secret value must not be empty")

	// validKeyRe matches POSIX-compliant environment variable names.
	validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

// EnvKey validates and normalises an environment variable name.
// It trims surrounding whitespace, upper-cases the result, and
// returns ErrInvalidEnvKey if the name does not match the POSIX
// naming rules after normalisation.
func EnvKey(raw string) (string, error) {
	key := strings.TrimSpace(raw)
	if key == "" {
		return "", ErrInvalidEnvKey
	}
	key = strings.ToUpper(key)
	if !validKeyRe.MatchString(key) {
		return "", ErrInvalidEnvKey
	}
	return key, nil
}

// SecretValue trims surrounding whitespace from a secret string and
// returns ErrEmptyValue when the result is empty.
func SecretValue(raw string) (string, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return "", ErrEmptyValue
	}
	return v, nil
}

// EnvPair sanitises both the key and value of a prospective environment
// entry, returning the cleaned pair or the first error encountered.
func EnvPair(key, value string) (string, string, error) {
	cleanKey, err := EnvKey(key)
	if err != nil {
		return "", "", err
	}
	cleanVal, err := SecretValue(value)
	if err != nil {
		return "", "", err
	}
	return cleanKey, cleanVal, nil
}
