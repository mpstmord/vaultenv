// Package redact provides a Redactor type that scrubs known secret values
// from plain strings and environment maps, replacing them with [REDACTED].
// It is intended to be used as a final safety layer before any output is
// written to logs, audit trails, or terminal output.
package redact
