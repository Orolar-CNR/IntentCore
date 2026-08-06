package state

import (
	"context"
	"errors"
	"sync"

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
	mu     sync.RWMutex
	latest *InternalSnapshot
	data   map[string]*contracts.Snapshot
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
	return &InMemorySnapshotStore{
		data: make(map[string]*contracts.Snapshot),
	}
}

func (s *InMemorySnapshotStore) Save(ctx context.Context, snapshot *contracts.Snapshot) error {
	if snapshot == nil {
		return errors.New("cannot save nil snapshot")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cp := *snapshot
	s.data[snapshot.ID] = &cp
	s.latest = &InternalSnapshot{
		Header: cp,
		Data:   make(map[string]*stateEntry),
	}
	return nil
}

func (s *InMemorySnapshotStore) SaveInternal(ctx context.Context, snapshot any) error {
	snap, ok := snapshot.(InternalSnapshot)
	if !ok {
		return errors.New("invalid snapshot type for in-memory store")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.latest = &snap
	cp := snap.Header
	s.data[cp.ID] = &cp
	return nil
}

func (s *InMemorySnapshotStore) Load(ctx context.Context, id string) (*contracts.Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, ok := s.data[id]
	if !ok {
		return nil, nil
	}

	cp := *snapshot
	return &cp, nil
}

func (s *InMemorySnapshotStore) LoadLatest(ctx context.Context) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.latest == nil {
		// Return an empty snapshot if none exists
		return InternalSnapshot{
			Header: contracts.Snapshot{
				ID:          uuid.New().String(),
				Offset:      0,
				IntentCount: 0,
			},
			Data: make(map[string]*stateEntry),
		}, nil
	}
	return *s.latest, nil
}

func (s *InMemorySnapshotStore) List(ctx context.Context) ([]contracts.Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]contracts.Snapshot, 0, len(s.data))
	for _, snapshot := range s.data {
		out = append(out, *snapshot)
	}
	return out, nil
}

func (s *InMemorySnapshotStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
