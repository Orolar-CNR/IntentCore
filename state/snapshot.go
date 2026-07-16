package state

import (
	"context"
	"errors"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/google/uuid"
)

// InternalSnapshot represents an internal, detailed checkpoint containing actual state entries.
type InternalSnapshot struct {
	Header contracts.Snapshot
	Data   map[string]*stateEntry // string representation of IntentID
}

// InMemorySnapshotStore implements contracts.SnapshotStore for testing.
type InMemorySnapshotStore struct {
	latest *InternalSnapshot
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
	return &InMemorySnapshotStore{}
}

func (s *InMemorySnapshotStore) Save(ctx context.Context, snapshot any) error {
	snap, ok := snapshot.(InternalSnapshot)
	if !ok {
		return errors.New("invalid snapshot type for in-memory store")
	}
	s.latest = &snap
	return nil
}

func (s *InMemorySnapshotStore) LoadLatest(ctx context.Context) (any, error) {
	if s.latest == nil {
		// Return an empty snapshot if none exists
		return InternalSnapshot{
			Header: contracts.Snapshot{
				SnapshotID:  uuid.New().String(),
				Offset:      0,
				IntentCount: 0,
			},
			Data: make(map[string]*stateEntry),
		}, nil
	}
	return *s.latest, nil
}
