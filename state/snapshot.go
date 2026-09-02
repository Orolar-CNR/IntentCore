package state

import (
	"context"
	"errors"
	"sync"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// InternalSnapshot is the full snapshot including headers and data.
type InternalSnapshot struct {
	Header contracts.Snapshot
	Data   map[string]*stateEntry
}

type InMemorySnapshotStore struct {
	mu     sync.RWMutex
	latest *InternalSnapshot
	data   map[string]*contracts.Snapshot // headers, keyed by snapshot ID

	// internalData holds the full InternalSnapshot (including state
	// entries) keyed by snapshot ID, so a specific historical snapshot
	// can be recovered by ID rather than only the most recent one.
	internalData map[string]*InternalSnapshot
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
	return &InMemorySnapshotStore{
		data:         make(map[string]*contracts.Snapshot),
		internalData: make(map[string]*InternalSnapshot),
	}
}

// Ensure InMemorySnapshotStore implements SnapshotStoreExtended.
var _ SnapshotStoreExtended = (*InMemorySnapshotStore)(nil)

func (s *InMemorySnapshotStore) Save(ctx context.Context, snapshot *contracts.Snapshot) error {
	if snapshot == nil {
		return errors.New("cannot save nil snapshot")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	cp := *snapshot
	s.data[snapshot.ID] = &cp
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

	snapCopy := snap
	s.internalData[cp.ID] = &snapCopy
	return nil
}

// LoadInternal returns the full InternalSnapshot (including state entries)
// for the given snapshot ID, or (nil, nil) if no such snapshot exists.
func (s *InMemorySnapshotStore) LoadInternal(ctx context.Context, id string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap, ok := s.internalData[id]
	if !ok {
		return nil, nil
	}
	cp := *snap
	return cp, nil
}

func (s *InMemorySnapshotStore) Load(ctx context.Context, id string) (*contracts.Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	cp := *snap
	return &cp, nil
}

func (s *InMemorySnapshotStore) LoadLatest(ctx context.Context) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.latest == nil {
		return nil, nil
	}
	cp := *s.latest
	return cp, nil
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
	delete(s.internalData, id)
	return nil
}
