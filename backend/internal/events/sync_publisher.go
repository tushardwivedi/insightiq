package events

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// SyncEventPublisher is a synchronous implementation of EventPublisher
// It publishes events immediately to all registered handlers
type SyncEventPublisher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
	logger   *slog.Logger
}

// NewSyncEventPublisher creates a new synchronous event publisher
func NewSyncEventPublisher(logger *slog.Logger) *SyncEventPublisher {
	return &SyncEventPublisher{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

// Subscribe registers a handler for a specific event type
func (p *SyncEventPublisher) Subscribe(eventType string, handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.handlers[eventType] == nil {
		p.handlers[eventType] = make([]EventHandler, 0)
	}

	p.handlers[eventType] = append(p.handlers[eventType], handler)

	p.logger.Info("Event handler subscribed",
		"event_type", eventType,
		"total_handlers", len(p.handlers[eventType]))
}

// Publish publishes an event to all registered handlers for the given event type
func (p *SyncEventPublisher) Publish(ctx context.Context, eventType string, event interface{}) error {
	p.mu.RLock()
	handlers := p.handlers[eventType]
	p.mu.RUnlock()

	if len(handlers) == 0 {
		p.logger.Debug("No handlers registered for event type",
			"event_type", eventType)
		return nil
	}

	p.logger.Debug("Publishing event",
		"event_type", eventType,
		"handler_count", len(handlers))

	var lastErr error
	for i, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			// Log the error but continue processing other handlers
			// This ensures one failing handler doesn't break the entire event processing
			p.logger.Error("Event handler failed",
				"event_type", eventType,
				"handler_index", i,
				"error", err)
			lastErr = err
		}
	}

	if lastErr != nil {
		return fmt.Errorf("one or more event handlers failed: %w", lastErr)
	}

	return nil
}
