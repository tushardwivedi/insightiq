package events

import (
	"context"
)

// EventPublisher defines the interface for publishing events
// This abstraction allows for different implementations (sync, async, message queue, etc.)
type EventPublisher interface {
	// Publish publishes an event to all registered handlers
	Publish(ctx context.Context, eventType string, event interface{}) error

	// Subscribe registers a handler for a specific event type
	Subscribe(eventType string, handler EventHandler)
}

// EventHandler defines the interface for handling events
type EventHandler interface {
	// Handle processes an event
	Handle(ctx context.Context, event interface{}) error
}
