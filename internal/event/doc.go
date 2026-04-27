// Package event implements a lightweight synchronous pub/sub event bus used
// by vaultenv components to broadcast lifecycle events such as secret
// rotation, token renewal, and health state changes.
//
// Subscribers register a Handler for a specific event.Type; the Bus delivers
// published events to all matching handlers in registration order.
package event
