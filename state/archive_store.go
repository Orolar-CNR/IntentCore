package state

import (
	"context"
	"sync"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// InMemoryArchiveStore is a cold-storage baseline implementation.
type InMemoryArchiveStore struct {
	mu   sync.RWMutex
	data map[string]*contracts.ArchivedSnapshot
}

func NewInMemoryArchiveStore() *InMemoryArchiveStore {
	return &InMemoryArchiveStore{
		data: make(map[string]*contracts.ArchivedSnapshot),
	}
}

func (s *InMemoryArchiveStore) Archive(ctx context.Context, entry *contracts.ArchivedSnapshot) error {
	if entry == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cp := *entry
	s.data[entry.ID] = &cp
	return nil
}

func (s *InMemoryArchiveStore) Retrieve(ctx context.Context, id string) (*contracts.ArchivedSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.data[id]
	if !ok {
		return nil, nil
	}

	cp := *entry
	return &cp, nil
}

func (s *InMemoryArchiveStore) List(ctx context.Context) ([]contracts.ArchivedSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]contracts.ArchivedSnapshot, 0, len(s.data))
	for _, entry := range s.data {
		out = append(out, *entry)
	}
	return out, nil
}

func (s *InMemoryArchiveStore) Purge(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}
