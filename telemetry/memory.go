package telemetry

import (
	"context"
	"sync"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// InMemoryRecorder implements contracts.TelemetryRecorder for testing and local usage.
type InMemoryRecorder struct {
	mu     sync.RWMutex
	events []contracts.TelemetryEvent
}

// NewInMemoryRecorder creates a new local telemetry recorder.
func NewInMemoryRecorder() *InMemoryRecorder {
	return &InMemoryRecorder{
		events: make([]contracts.TelemetryEvent, 0),
	}
}

// Record appends a telemetry event to the in-memory buffer.
func (r *InMemoryRecorder) Record(ctx context.Context, event contracts.TelemetryEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}

// GetEvents returns a copy of the recorded events for inspection.
func (r *InMemoryRecorder) GetEvents() []contracts.TelemetryEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cpy := make([]contracts.TelemetryEvent, len(r.events))
	copy(cpy, r.events)
	return cpy
}
