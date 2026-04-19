// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound requests to the HashiCorp Vault API, preventing accidental
// exhaustion of server-side rate limits.
package ratelimit
