package contracts

import (
	"context"
	"time"

	"github.com/Orolar-CNR/IntentCore/core"
)

// Snapshot represents a checkpoint of the repository state at a specific ledger offset.
type Snapshot struct {
	ID            string
	SchemaVersion core.StateVersion
	CreatedAt     time.Time
	Checkpoint    string
	IntentCount   uint64
	Payload       []byte

	// Checksum is the hex-encoded SHA-256 digest of the snapshot's state
	// data at the time it was created by Repository.Snapshot. Recover
	// MUST recompute this digest from the loaded data and reject the
	// snapshot on mismatch, to detect tampering or corruption of
	// snapshots at rest (SnapshotStore / ArchiveStore).
	Checksum string

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
