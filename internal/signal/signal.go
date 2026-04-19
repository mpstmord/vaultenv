// Package signal handles OS signal forwarding to child processes.
package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Forwarder listens for OS signals and forwards them to a target process.
type Forwarder struct {
	proc   *os.Process
	signals chan os.Signal
	done   chan struct{}
}

// NewForwarder creates a Forwarder that will forward signals to proc.
func NewForwarder(proc *os.Process) *Forwarder {
	f := &Forwarder{
		proc:   proc,
		signals: make(chan os.Signal, 8),
		done:   make(chan struct{}),
	}
	return f
}

// Start begins forwarding signals. Call Stop to clean up.
func (f *Forwarder) Start() {
	signal.Notify(f.signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	go func() {
		for {
			select {
			case sig := <-f.signals:
				_ = f.proc.Signal(sig)
			case <-f.done:
				return
			}
		}
	}()
}

// Stop stops forwarding signals and unregisters the notify channel.
func (f *Forwarder) Stop() {
	signal.Stop(f.signals)
	close(f.done)
}
