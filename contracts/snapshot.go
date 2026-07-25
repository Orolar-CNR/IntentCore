package contracts

import (
	"context"
)

// SnapshotStore defines persistence for repository snapshots.
type SnapshotStore interface {
	Save(ctx context.Context, snapshot *Snapshot) error
	Load(ctx context.Context, id string) (*Snapshot, error)
	List(ctx context.Context) ([]Snapshot, error)
	Delete(ctx context.Context, id string) error
}
