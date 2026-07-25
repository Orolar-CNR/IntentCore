package state

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

func TestDefaultSnapshotScheduler(t *testing.T) {
	s := DefaultSnapshotScheduler{IntervalCount: 10}
	ctx := context.Background()
	now := time.Now()

	if s.ShouldSnapshot(ctx, now, 0) {
		t.Error("Should not snapshot on 0 count")
	}

	if s.ShouldSnapshot(ctx, now, 5) {
		t.Error("Should not snapshot on 5 count")
	}

	if !s.ShouldSnapshot(ctx, now, 10) {
		t.Error("Should snapshot on 10 count")
	}

	if !s.ShouldSnapshot(ctx, now, 20) {
		t.Error("Should snapshot on 20 count")
	}
}

func TestDefaultRetentionPolicy(t *testing.T) {
	p := DefaultRetentionPolicy{ArchiveAfterSnapshots: 100}
	ctx := context.Background()

	decision, err := p.Decide(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if decision != contracts.RetentionDelete {
		t.Errorf("Expected RetentionDelete for nil snapshot, got %v", decision)
	}

	snapshot := &contracts.Snapshot{IntentCount: 50}
	decision, err = p.Decide(ctx, snapshot)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if decision != contracts.RetentionKeep {
		t.Errorf("Expected RetentionKeep for intent count 50, got %v", decision)
	}

	snapshot.IntentCount = 100
	decision, err = p.Decide(ctx, snapshot)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if decision != contracts.RetentionArchive {
		t.Errorf("Expected RetentionArchive for intent count 100, got %v", decision)
	}
}
