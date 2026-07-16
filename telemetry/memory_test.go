package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/telemetry"
)

func TestInMemoryRecorder(t *testing.T) {
	recorder := telemetry.NewInMemoryRecorder()

	event := contracts.TelemetryEvent{
		IntentID:  "test-intent",
		EventType: "Admission_Accepted",
		Timestamp: time.Now(),
	}

	err := recorder.Record(context.Background(), event)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	events := recorder.GetEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].IntentID != "test-intent" {
		t.Errorf("Expected IntentID test-intent, got %s", events[0].IntentID)
	}
}
