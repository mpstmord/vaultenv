// Package event provides a simple pub/sub event bus for broadcasting
// internal lifecycle events across vaultenv components.
package event

import (
	"sync"
)

// Type identifies the kind of event.
type Type string

const (
	SecretFetched  Type = "secret.fetched"
	SecretRotated  Type = "secret.rotated"
	SecretExpired  Type = "secret.expired"
	TokenRenewed   Type = "token.renewed"
	HealthChanged  Type = "health.changed"
)

// Event carries a type and an arbitrary payload.
type Event struct {
	Type    Type
	Payload any
}

// Handler is a function that receives an event.
type Handler func(Event)

// Bus is a simple synchronous event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers map[Type][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[Type][]Handler),
	}
}

// Subscribe registers h to receive events of the given type.
// Calling Subscribe with the same handler multiple times registers it
// multiple times; callers are responsible for deduplication if needed.
func (b *Bus) Subscribe(t Type, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

// Publish delivers e to all handlers registered for e.Type.
// Handlers are called synchronously in registration order.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers[e.Type]))
	copy(handlers, b.handlers[e.Type])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

// Unsubscribe removes all handlers registered for the given type.
func (b *Bus) Unsubscribe(t Type) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, t)
}
