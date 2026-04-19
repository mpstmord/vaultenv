// Package signal provides OS signal forwarding for child processes spawned
// by vaultenv, ensuring that signals such as SIGTERM and SIGHUP are
// propagated correctly so secrets-injected processes shut down gracefully.
package signal
