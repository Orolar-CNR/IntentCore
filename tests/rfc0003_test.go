package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/google/uuid"
)

func TestRFC0003_CAS(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo := state.NewRepository(store)

	id := core.IntentID(uuid.New())

	// 1. Initial State Transition
	record1 := contracts.IntentRecord{
		ID:        id,
		Agent:     "test-actor",
		Timestamp: time.Now(),
		State:     contracts.StatePending,
		Version:   1, // Expected to be 1 after CAS of 0
	}

	err := repo.CompareAndSwap(ctx, 0, record1)
	if err != nil {
		t.Fatalf("Expected no error on initial CAS, got %v", err)
	}

	// 2. Load Intent Verification
	loaded, err := repo.LoadIntent(ctx, id)
	if err != nil {
		t.Fatalf("Expected no error loading intent, got %v", err)
	}
	if loaded.Version != 1 {
		t.Errorf("Expected version 1, got %d", loaded.Version)
	}

	// 3. Version Conflict Detection
	record2 := contracts.IntentRecord{
		ID:        id,
		Agent:     "test-actor",
		Timestamp: time.Now(),
		State:     contracts.StateExecuting,
		Version:   2,
	}
	
	// Intentionally provide the wrong expected version
	err = repo.CompareAndSwap(ctx, 5, record2)
	if err == nil {
		t.Fatalf("Expected version conflict error, got nil")
	}

	// 4. Successful CAS Progression
	err = repo.CompareAndSwap(ctx, 1, record2)
	if err != nil {
		t.Fatalf("Expected no error on valid CAS, got %v", err)
	}
}
