package contracts

import (
	"context"
)

// Snapshot represents a checkpoint of the repository state at a specific ledger offset.
type Snapshot struct {
	SnapshotID  string
	Offset      uint64
	IntentCount int64
	// In a real implementation, this would contain a stream or reference to the blob data.
}

// StateRepository defines the canonical state persistence contract.
//
// RFC:
//
//	RFC-0003 Section 4
//
// Guarantees:
//   - Single Source of Truth
//   - Compare-and-Swap
//   - Immutable History
type StateRepository interface {
	// LoadIntent retrieves the current state and version of an Intent.
	LoadIntent(ctx context.Context, id IntentID) (*IntentRecord, error)

	// CompareAndSwap atomically updates the state of an Intent if the expected version matches.
	// Returns core.ErrVersionConflict if the version does not match.
	CompareAndSwap(ctx context.Context, expected Version, next IntentRecord) error

	// Snapshot creates a durable checkpoint of the repository state.
	Snapshot(ctx context.Context) (*Snapshot, error)

	// Recover restores the repository state from a snapshot and replays subsequent ledger events.
	Recover(ctx context.Context, snapshot Snapshot) error
}
