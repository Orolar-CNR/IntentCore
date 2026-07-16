package state_test

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/google/uuid"
)

func TestRepository_SnapshotAndRecover(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo1 := state.NewRepository(store)

	// Add an intent
	id := core.IntentID(uuid.New())
	record := contracts.IntentRecord{
		ID:        id,
		Agent:     "test-agent",
		State:     contracts.StatePending,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   []byte("test"),
	}

	err := repo1.CompareAndSwap(ctx, 0, record)
	if err != nil {
		t.Fatalf("Failed to CAS: %v", err)
	}

	// Take Snapshot
	header, err := repo1.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Failed to snapshot: %v", err)
	}

	if header.IntentCount != 1 {
		t.Errorf("Expected 1 intent in snapshot header, got %d", header.IntentCount)
	}

	// Create a new Repository and recover from the snapshot
	repo2 := state.NewRepository(store)
	
	err = repo2.Recover(ctx, *header)
	if err != nil {
		t.Fatalf("Failed to recover: %v", err)
	}

	// Verify the data exists in the new repo
	loaded, err := repo2.LoadIntent(ctx, id)
	if err != nil {
		t.Fatalf("Failed to load intent after recovery: %v", err)
	}

	if loaded.State != contracts.StatePending {
		t.Errorf("Expected state Pending, got %s", loaded.State)
	}
	if loaded.Version != 1 {
		t.Errorf("Expected version 1, got %d", loaded.Version)
	}
}
