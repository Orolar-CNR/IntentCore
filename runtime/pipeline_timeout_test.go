package runtime_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/transport"
	"github.com/google/uuid"
)

// slowEvaluator mimics a policy that takes a long time to evaluate
type slowEvaluator struct {
	delay time.Duration
}

func (s *slowEvaluator) Evaluate(ctx context.Context, env contracts.SemanticEnvelope) (contracts.AdmissionResult, error) {
	select {
	case <-time.After(s.delay):
		return contracts.AdmissionResult{
			Decision: contracts.DecisionAccept,
		}, nil
	case <-ctx.Done():
		return contracts.AdmissionResult{}, ctx.Err()
	}
}

func TestPipelineTimeoutAndCancellation(t *testing.T) {
	// Create a payload
	env := contracts.SemanticEnvelope{
		SchemaVersion:  admission.DefaultSchemaVersion,
		EnvelopeID:     core.IntentID(uuid.New()),
		AgentIdentity:  "test-agent",
		EventTimestamp: time.Now(),
		Signatures:     []contracts.Signature{{Signature: []byte("test-sig")}},
		OpaquePayload:  []byte(`{"key":"value"}`),
	}
	payloadBytes, _ := json.Marshal(env)

	// Setup slow evaluator that takes 500ms
	evaluator := &slowEvaluator{delay: 500 * time.Millisecond}
	trans := transport.NewMockTransport([][]byte{payloadBytes})

	// Create pipeline with 100ms timeout
	pipeline := runtime.NewPipeline(
		trans,
		evaluator,
		&mockLifecycle{}, // from dispatcher_test.go
		runtime.WithTimeout(100*time.Millisecond),
		runtime.WithRetries(0),
	)

	// Start pipeline, dispatcher worker will sit idle
	pipeline.Start(context.Background())
	defer pipeline.Stop()

	// Execute should timeout before slowEvaluator finishes
	err := pipeline.Execute(context.Background(), payloadBytes)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected DeadlineExceeded, got: %v", err)
	}
}
