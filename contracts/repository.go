package contracts

import (
	"context"
	"github.com/Orolar-CNR/IntentCore/core"
	"time"
)

// Snapshot represents a checkpoint of the repository state at a specific ledger offset.
type Snapshot struct {
	ID            string
	SchemaVersion core.StateVersion
	CreatedAt     time.Time
	Checkpoint    string
	IntentCount   uint64
	Payload       []byte

	SnapshotID string // Deprecated, use ID
	Offset     uint64
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
	LoadIntent(ctx context.Context, id core.IntentID) (*IntentRecord, error)

	// CompareAndSwap atomically updates the state of an Intent if the expected version matches.
	// Returns core.ErrVersionConflict if the version does not match.
	CompareAndSwap(ctx context.Context, expected core.StateVersion, next IntentRecord) error

	// Snapshot creates a durable checkpoint of the repository state.
	Snapshot(ctx context.Context) (*Snapshot, error)

	// Recover restores the repository state from a snapshot and replays subsequent ledger events.
	Recover(ctx context.Context, snapshot Snapshot) error
}
