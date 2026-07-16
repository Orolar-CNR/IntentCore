package state

import (
	"context"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// SnapshotScheduler decides when snapshots should be produced.
type SnapshotScheduler interface {
	ShouldSnapshot(ctx context.Context, now time.Time, snapshotCount uint64) bool
}

// DefaultSnapshotScheduler is a simple threshold-based scheduler.
type DefaultSnapshotScheduler struct {
	IntervalCount uint64
}

// ShouldSnapshot returns true when the snapshot count reaches the configured interval.
func (s DefaultSnapshotScheduler) ShouldSnapshot(ctx context.Context, now time.Time, snapshotCount uint64) bool {
	if s.IntervalCount == 0 {
		return false
	}
	return snapshotCount > 0 && snapshotCount%s.IntervalCount == 0
}

var _ SnapshotScheduler = DefaultSnapshotScheduler{}
var _ contracts.RetentionPolicy = (*DefaultRetentionPolicy)(nil)

// DefaultRetentionPolicy is a simple policy for Phase 3.1.
type DefaultRetentionPolicy struct {
	ArchiveAfterSnapshots uint64
}

// Decide determines whether the snapshot should be kept or archived.
func (p *DefaultRetentionPolicy) Decide(ctx context.Context, snapshot *contracts.Snapshot) (contracts.RetentionDecision, error) {
	if snapshot == nil {
		return contracts.RetentionDelete, nil
	}

	if p.ArchiveAfterSnapshots > 0 && snapshot.IntentCount >= p.ArchiveAfterSnapshots {
		return contracts.RetentionArchive, nil
	}

	return contracts.RetentionKeep, nil
}
