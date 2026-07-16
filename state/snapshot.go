package state

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/google/uuid"
)

// Snapshot represents an internal, detailed checkpoint containing actual state entries.
type InternalSnapshot struct {
	Header contracts.Snapshot
	Data   map[string]*stateEntry // string representation of IntentID
}

// SnapshotStore defines the abstraction for saving and loading snapshots.
// This allows swapping between Memory, File, BoltDB, or S3.
type SnapshotStore interface {
	Save(ctx context.Context, snapshot InternalSnapshot) error
	LoadLatest(ctx context.Context) (*InternalSnapshot, error)
}

// InMemorySnapshotStore implements SnapshotStore for testing.
type InMemorySnapshotStore struct {
	latest *InternalSnapshot
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
	return &InMemorySnapshotStore{}
}

func (s *InMemorySnapshotStore) Save(ctx context.Context, snapshot InternalSnapshot) error {
	s.latest = &snapshot
	return nil
}

func (s *InMemorySnapshotStore) LoadLatest(ctx context.Context) (*InternalSnapshot, error) {
	if s.latest == nil {
		// Return an empty snapshot if none exists
		return &InternalSnapshot{
			Header: contracts.Snapshot{
				SnapshotID:  uuid.New().String(),
				Offset:      0,
				IntentCount: 0,
			},
			Data: make(map[string]*stateEntry),
		}, nil
	}
	return s.latest, nil
}
