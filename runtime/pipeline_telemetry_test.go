package runtime_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/telemetry"
	"github.com/Orolar-CNR/IntentCore/transport"
	"github.com/google/uuid"
)

func TestPipelineTelemetry(t *testing.T) {
	// Create a valid payload
	env := contracts.SemanticEnvelope{
		SchemaVersion:  admission.DefaultSchemaVersion,
		EnvelopeID:     core.IntentID(uuid.New()),
		AgentIdentity:  "test-agent",
		EventTimestamp: time.Now(),
		Signatures:     []contracts.Signature{{Signature: []byte("test-sig")}},
		OpaquePayload:  []byte(`{"key":"value"}`),
	}
	payloadBytes, _ := json.Marshal(env)

	evaluator := admission.NewPolicyEvaluator()
	trans := transport.NewMockTransport([][]byte{payloadBytes})

	recorder := telemetry.NewInMemoryRecorder()

	pipeline := runtime.NewPipeline(
		trans,
		evaluator,
		&mockLifecycle{}, // from dispatcher_test.go
		runtime.WithTelemetry(recorder),
	)

	// Execute should emit telemetry
	err := pipeline.Execute(context.Background(), payloadBytes)
	if err != nil {
		t.Fatalf("Pipeline execute failed: %v", err)
	}

	events := recorder.GetEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 telemetry event, got %d", len(events))
	}

	if events[0].EventType != "Admission_Accepted" {
		t.Errorf("Expected Admission_Accepted, got %s", events[0].EventType)
	}
}

func TestPipelineTelemetryReject(t *testing.T) {
	// Create an invalid payload (missing identity)
	env := contracts.SemanticEnvelope{
		SchemaVersion:  admission.DefaultSchemaVersion,
		EnvelopeID:     core.IntentID(uuid.New()),
		EventTimestamp: time.Now(),
		Signatures:     []contracts.Signature{{Signature: []byte("test-sig")}},
		OpaquePayload:  []byte(`{"key":"value"}`),
	}
	payloadBytes, _ := json.Marshal(env)

	evaluator := admission.NewPolicyEvaluator()
	trans := transport.NewMockTransport([][]byte{payloadBytes})

	recorder := telemetry.NewInMemoryRecorder()

	pipeline := runtime.NewPipeline(
		trans,
		evaluator,
		&mockLifecycle{},
		runtime.WithTelemetry(recorder),
	)

	// Execute should fail and emit reject telemetry
	err := pipeline.Execute(context.Background(), payloadBytes)
	if err == nil {
		t.Fatal("Expected pipeline execute to fail")
	}

	events := recorder.GetEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 telemetry event, got %d", len(events))
	}

	if events[0].EventType != "Admission_Rejected" {
		t.Errorf("Expected Admission_Rejected, got %s", events[0].EventType)
	}
}
