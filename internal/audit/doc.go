// Package audit provides structured, append-only audit logging for vaultenv.
//
// Each secret fetch, field resolution, and process exec is recorded as a
// newline-delimited JSON event. Logging is optional; pass nil to NewLogger
// to disable it entirely with zero overhead.
package audit
