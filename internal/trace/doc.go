// Package trace provides lightweight span-based tracing for vaultenv.
//
// A Span is created at the start of each secret fetch operation and
// stored in the request context. Middleware and handlers can retrieve
// the active span via FromContext to attach trace IDs to audit log
// entries, enabling end-to-end correlation across the fetch pipeline.
package trace
