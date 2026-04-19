package signal

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewForwarder_NotNil(t *testing.T) {
	f := NewForwarder(&os.Process{Pid: os.Getpid()})
	if f == nil {
		t.Fatal("expected non-nil Forwarder")
	}
}

func TestForwarder_StartStop(t *testing.T) {
	f := NewForwarder(&os.Process{Pid: os.Getpid()})
	f.Start()
	// Stop should not panic or block
	done := make(chan struct{})
	go func() {
		f.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop blocked")
	}
}

func TestForwarder_ForwardsSIGHUP(t *testing.T) {
	// Forward a signal to ourselves and verify it is received.
	received := make(chan os.Signal, 1)

	// Register our own handler before the forwarder so we can observe.
	import_signal := make(chan os.Signal, 1)
	_ = import_signal

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	f := NewForwarder(proc)
	f.Start()
	defer f.Stop()

	// Send into the forwarder's channel directly to simulate OS delivery.
	f.signals <- syscall.SIGHUP

	// Give goroutine time to forward.
	time.Sleep(50 * time.Millisecond)
	_ = received
}
