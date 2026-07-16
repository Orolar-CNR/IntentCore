package contracts

import "context"

// SnapshotStore defines the abstraction for saving and loading snapshots.
// This allows swapping between Memory, File, BoltDB, or S3.
type SnapshotStore interface {
	Save(ctx context.Context, snapshot any) error
	LoadLatest(ctx context.Context) (any, error)
}
