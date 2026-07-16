package state

import (
	"context"
	"sync"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// Repository implements contracts.StateRepository.
// It maintains the single source of truth for the system, enforcing CAS semantics.
type Repository struct {
	mu    sync.RWMutex
	store map[core.IntentID]*stateEntry
}

// NewRepository initializes a new State Repository.
func NewRepository() *Repository {
	return &Repository{
		store: make(map[core.IntentID]*stateEntry),
	}
}

func (r *Repository) LoadIntent(ctx context.Context, id core.IntentID) (*contracts.IntentRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.store[id]
	if !exists {
		return nil, core.ErrNotFound
	}

	return entry.toRecord(), nil
}

func (r *Repository) CompareAndSwap(ctx context.Context, expected core.StateVersion, next contracts.IntentRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, exists := r.store[next.ID]
	var currentVersion core.StateVersion = 0
	if exists {
		currentVersion = entry.Version
	}

	if err := EvaluateCAS(expected, currentVersion); err != nil {
		return err
	}

	// Update or insert the new entry
	r.store[next.ID] = &stateEntry{
		ID:        next.ID,
		Agent:     next.Agent,
		Timestamp: next.Timestamp,
		Payload:   next.Payload,
		State:     next.State,
		Version:   next.Version, // Should be expected + 1, managed by lifecycle
	}

	return nil
}

func (r *Repository) Snapshot(ctx context.Context) (*contracts.Snapshot, error) {
	panic("not implemented: see RFC-0003")
}

func (r *Repository) Recover(ctx context.Context, snapshot contracts.Snapshot) error {
	panic("not implemented: see RFC-0003")
}
