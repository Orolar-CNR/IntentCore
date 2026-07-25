package runtime_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/proof"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/transport"
	"github.com/google/uuid"
)

func TestPipelineProof(t *testing.T) {
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

	recorder := proof.NewInMemoryProofRecorder()

	pipeline := runtime.NewPipeline(
		trans,
		evaluator,
		&mockLifecycle{}, // from dispatcher_test.go
		runtime.WithProof(recorder),
	)

	// Execute should emit proof
	err := pipeline.Execute(context.Background(), payloadBytes)
	if err != nil {
		t.Fatalf("Pipeline execute failed: %v", err)
	}

	proofs := recorder.GetProofs()
	if len(proofs) != 1 {
		t.Fatalf("Expected 1 proof event, got %d", len(proofs))
	}

	if proofs[0].IntentID != uuid.UUID(env.EnvelopeID).String() {
		t.Errorf("Expected Proof for %s, got %s", env.EnvelopeID, proofs[0].IntentID)
	}
}
