// Package circuit provides a circuit breaker for protecting upstream secret
// fetches from cascading failures. When a configurable failure threshold is
// reached the breaker opens and immediately rejects requests with ErrOpen.
// After a reset timeout it transitions to half-open, allowing a probe request
// through; a success closes the circuit while another failure re-opens it.
package circuit
