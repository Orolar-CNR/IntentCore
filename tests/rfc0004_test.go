package tests

import (
	"context"
	"testing"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/history"
	"github.com/Orolar-CNR/IntentCore/lifecycle"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/google/uuid"
)

func TestRFC0004_TransitionMatrix(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo := state.NewRepository(store)
	ledger := history.NewLedger()
	recorder := history.NewRecorder(ledger)
	machine := lifecycle.NewStateMachine(repo, recorder)

	id := core.IntentID(uuid.New())

	// Test 1: Valid Transition (Initial -> Pending)
	req1 := contracts.TransitionRequest{
		IntentID:        id,
		FromState:       "", // initial
		ToState:         contracts.StatePending,
		ExpectedVersion: 0,
		ActorID:         "test-actor",
		Metadata: map[string]any{
			"Payload": []byte("test"),
		},
	}
	err := machine.Transition(ctx, req1)
	if err != nil {
		t.Fatalf("Expected no error on valid transition to Pending, got %v", err)
	}

	// Test 2: Valid Transition (Pending -> Validated)
	req2 := contracts.TransitionRequest{
		IntentID:        id,
		FromState:       contracts.StatePending,
		ToState:         contracts.StateValidated,
		ExpectedVersion: 1,
		ActorID:         "test-actor",
	}
	err = machine.Transition(ctx, req2)
	if err != nil {
		t.Fatalf("Expected no error on valid transition to Validated, got %v", err)
	}

	// Test 3: Invalid Transition (Validated -> Pending)
	req3 := contracts.TransitionRequest{
		IntentID:        id,
		FromState:       contracts.StateValidated,
		ToState:         contracts.StatePending, // Backwards transition is not allowed
		ExpectedVersion: 2,
		ActorID:         "test-actor",
	}
	err = machine.Transition(ctx, req3)
	if err == nil {
		t.Fatalf("Expected ErrInvalidTransition, got nil")
	}
	
	// Test 4: Walk to Executing -> Completed to test terminal states
	
	machine.Transition(ctx, contracts.TransitionRequest{IntentID: id, FromState: contracts.StateValidated, ToState: contracts.StateAdmitted, ExpectedVersion: 2, ActorID: "a"})
	machine.Transition(ctx, contracts.TransitionRequest{IntentID: id, FromState: contracts.StateAdmitted, ToState: contracts.StateScheduled, ExpectedVersion: 3, ActorID: "a"})
	machine.Transition(ctx, contracts.TransitionRequest{IntentID: id, FromState: contracts.StateScheduled, ToState: contracts.StateExecuting, ExpectedVersion: 4, ActorID: "a"})
	
	req4 := contracts.TransitionRequest{
		IntentID:        id,
		FromState:       contracts.StateExecuting,
		ToState:         contracts.StateCompleted,
		ExpectedVersion: 5,
		ActorID:         "test-actor",
	}
	err = machine.Transition(ctx, req4)
	if err != nil {
		t.Fatalf("Expected no error transitioning to terminal state, got %v", err)
	}

	// Test 5: Transition from Terminal State (Completed -> Failed)
	req5 := contracts.TransitionRequest{
		IntentID:        id,
		FromState:       contracts.StateCompleted,
		ToState:         contracts.StateFailed,
		ExpectedVersion: 6,
		ActorID:         "test-actor",
	}
	err = machine.Transition(ctx, req5)
	if err == nil {
		t.Fatalf("Expected error transitioning from terminal state, got nil")
	}
}
