package contracts

import (
	"context"
	"time"
)

// TelemetryEvent defines an observability event within the IntentCore kernel.
type TelemetryEvent struct {
	IntentID  string
	EventType string
	Timestamp time.Time
	Duration  time.Duration
	Metadata  map[string]any
}

// TelemetryRecorder provides an abstraction for observability and metrics gathering.
// This is strictly an OPTIONAL dependency for the Runtime Pipeline.
type TelemetryRecorder interface {
	Record(ctx context.Context, event TelemetryEvent) error
}
