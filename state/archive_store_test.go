package state

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

func TestInMemoryArchiveStore(t *testing.T) {
	store := NewInMemoryArchiveStore()
	ctx := context.Background()

	entry := &contracts.ArchivedSnapshot{
		ID:        "arch-1",
		CreatedAt: time.Now(),
	}

	err := store.Archive(ctx, entry)
	if err != nil {
		t.Fatalf("Failed to archive entry: %v", err)
	}

	loaded, err := store.Retrieve(ctx, "arch-1")
	if err != nil {
		t.Fatalf("Failed to retrieve entry: %v", err)
	}
	if loaded == nil || loaded.ID != "arch-1" {
		t.Errorf("Retrieved entry does not match")
	}

	list, err := store.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list entries: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 entry in list, got %d", len(list))
	}

	err = store.Purge(ctx, "arch-1")
	if err != nil {
		t.Fatalf("Failed to purge entry: %v", err)
	}

	loaded, err = store.Retrieve(ctx, "arch-1")
	if err != nil {
		t.Fatalf("Failed to retrieve entry after purge: %v", err)
	}
	if loaded != nil {
		t.Errorf("Expected nil entry after purge, got %v", loaded)
	}
}
