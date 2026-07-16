package runtime_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/history"
	"github.com/Orolar-CNR/IntentCore/lifecycle"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/Orolar-CNR/IntentCore/transport"
	"github.com/google/uuid"
)

func TestVerticalSlice(t *testing.T) {
	// Prepare test data
	intentID := core.IntentID(uuid.New())
	env := contracts.SemanticEnvelope{
		EnvelopeID:     intentID,
		AgentIdentity:  "test-agent",
		EventTimestamp: time.Now(),
		OpaquePayload:  []byte(`{"key":"value"}`),
	}
	payloadBytes, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// 1. Setup DI
	repo := state.NewRepository()
	ledger := history.NewLedger()
	recorder := history.NewRecorder(ledger)
	machine := lifecycle.NewStateMachine(repo, recorder)
	evaluator := admission.NewPolicyEvaluator()
	trans := transport.NewMockTransport([][]byte{payloadBytes})
	pipeline := runtime.NewPipeline(trans, evaluator, machine)

	// 2. Execute Pipeline
	ctx := context.Background()
	if err := pipeline.Start(ctx); err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}

	// 3. Verify Repository State
	record, err := repo.LoadIntent(ctx, intentID)
	if err != nil {
		t.Fatalf("Failed to load intent from repo: %v", err)
	}

	if record.State != contracts.StatePending {
		t.Errorf("Expected state %s, got %s", contracts.StatePending, record.State)
	}
	if record.Agent != "test-agent" {
		t.Errorf("Expected agent test-agent, got %s", record.Agent)
	}
	if record.Version != 1 {
		t.Errorf("Expected version 1, got %d", record.Version)
	}

	// 4. Verify History Ledger
	records := ledger.GetRecords(intentID)
	if len(records) != 1 {
		t.Fatalf("Expected 1 history record, got %d", len(records))
	}

	histRec := records[0]
	if histRec.ToState != contracts.StatePending {
		t.Errorf("Expected history to state %s, got %s", contracts.StatePending, histRec.ToState)
	}
	if histRec.RecordedVersion != 1 {
		t.Errorf("Expected history recorded version 1, got %d", histRec.RecordedVersion)
	}
}
