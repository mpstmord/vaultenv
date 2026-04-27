package event_test

import (
	"sync"
	"testing"

	"github.com/nicholasgasior/vaultenv/internal/event"
)

func TestPublish_NoSubscribers(t *testing.T) {
	b := event.New()
	// Must not panic when no handlers are registered.
	b.Publish(event.Event{Type: event.SecretFetched, Payload: "ok"})
}

func TestSubscribe_ReceivesEvent(t *testing.T) {
	b := event.New()

	var got event.Event
	b.Subscribe(event.SecretFetched, func(e event.Event) {
		got = e
	})

	want := event.Event{Type: event.SecretFetched, Payload: "path/to/secret"}
	b.Publish(want)

	if got.Type != want.Type {
		t.Errorf("type: got %q, want %q", got.Type, want.Type)
	}
	if got.Payload != want.Payload {
		t.Errorf("payload: got %v, want %v", got.Payload, want.Payload)
	}
}

func TestSubscribe_MultipleHandlers_AllCalled(t *testing.T) {
	b := event.New()

	var mu sync.Mutex
	var calls int
	inc := func(_ event.Event) {
		mu.Lock()
		calls++
		mu.Unlock()
	}

	b.Subscribe(event.SecretRotated, inc)
	b.Subscribe(event.SecretRotated, inc)
	b.Subscribe(event.SecretRotated, inc)

	b.Publish(event.Event{Type: event.SecretRotated})

	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestSubscribe_DifferentTypes_Isolated(t *testing.T) {
	b := event.New()

	var fetchedCalls, expiredCalls int
	b.Subscribe(event.SecretFetched, func(_ event.Event) { fetchedCalls++ })
	b.Subscribe(event.SecretExpired, func(_ event.Event) { expiredCalls++ })

	b.Publish(event.Event{Type: event.SecretFetched})

	if fetchedCalls != 1 {
		t.Errorf("fetched: got %d, want 1", fetchedCalls)
	}
	if expiredCalls != 0 {
		t.Errorf("expired: got %d, want 0", expiredCalls)
	}
}

func TestUnsubscribe_RemovesHandlers(t *testing.T) {
	b := event.New()

	var calls int
	b.Subscribe(event.TokenRenewed, func(_ event.Event) { calls++ })

	b.Publish(event.Event{Type: event.TokenRenewed})
	b.Unsubscribe(event.TokenRenewed)
	b.Publish(event.Event{Type: event.TokenRenewed})

	if calls != 1 {
		t.Errorf("expected 1 call after unsubscribe, got %d", calls)
	}
}

func TestEventType_Constants(t *testing.T) {
	types := []event.Type{
		event.SecretFetched,
		event.SecretRotated,
		event.SecretExpired,
		event.TokenRenewed,
		event.HealthChanged,
	}
	seen := make(map[event.Type]struct{})
	for _, tt := range types {
		if _, dup := seen[tt]; dup {
			t.Errorf("duplicate event type constant: %q", tt)
		}
		seen[tt] = struct{}{}
	}
}
