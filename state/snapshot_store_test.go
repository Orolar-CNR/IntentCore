package state

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

func TestInMemorySnapshotStore_SaveLoad(t *testing.T) {
	store := NewInMemorySnapshotStore()
	ctx := context.Background()

	snapshot := &contracts.Snapshot{
		ID:        "snap-1",
		CreatedAt: time.Now(),
	}

	err := store.Save(ctx, snapshot)
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	loaded, err := store.Load(ctx, "snap-1")
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}
	if loaded == nil || loaded.ID != "snap-1" {
		t.Errorf("Loaded snapshot does not match")
	}

	list, err := store.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 snapshot in list, got %d", len(list))
	}

	err = store.Delete(ctx, "snap-1")
	if err != nil {
		t.Fatalf("Failed to delete snapshot: %v", err)
	}

	loaded, err = store.Load(ctx, "snap-1")
	if err != nil {
		t.Fatalf("Failed to load snapshot after delete: %v", err)
	}
	if loaded != nil {
		t.Errorf("Expected nil snapshot after delete, got %v", loaded)
	}
}
