package state

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/google/uuid"
)

// SnapshotStoreExtended is an internal interface that extends SnapshotStore with methods needed for repository.go
type SnapshotStoreExtended interface {
	contracts.SnapshotStore
	SaveInternal(ctx context.Context, snapshot any) error
	LoadLatest(ctx context.Context) (any, error)
}

// Repository implements contracts.StateRepository.
// It maintains the single source of truth for the system, enforcing CAS semantics.
type Repository struct {
	mu            sync.RWMutex
	store         map[core.IntentID]*stateEntry
	snapshotStore SnapshotStoreExtended
}

// NewRepository initializes a new State Repository with an optional SnapshotStore.
// If none is provided, it defaults to InMemorySnapshotStore.
func NewRepository(store contracts.SnapshotStore) *Repository {
	var extendedStore SnapshotStoreExtended
	if store == nil {
		extendedStore = NewInMemorySnapshotStore()
	} else if ext, ok := store.(SnapshotStoreExtended); ok {
		extendedStore = ext
	} else {
		panic("provided store must implement SnapshotStoreExtended")
	}

	return &Repository{
		store:         make(map[core.IntentID]*stateEntry),
		snapshotStore: extendedStore,
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := make(map[string]*stateEntry)
	for k, v := range r.store {
		// deep copy the entry to prevent mutations
		clone := *v
		data[uuid.UUID(k).String()] = &clone
	}

	internalSnap := InternalSnapshot{
		Header: contracts.Snapshot{
			ID:          uuid.New().String(),
			Offset:      uint64(time.Now().UnixNano()), // Simplified offset for phase 2
			IntentCount: uint64(len(data)),
			SnapshotID:  uuid.New().String(), // Keep for backwards compatibility with tests
		},
		Data: data,
	}

	if err := r.snapshotStore.SaveInternal(ctx, internalSnap); err != nil {
		return nil, err
	}

	return &internalSnap.Header, nil
}

func (r *Repository) Recover(ctx context.Context, snapshot contracts.Snapshot) error {
	// Re-hydrate state from the snapshot.
	loaded, err := r.snapshotStore.LoadLatest(ctx)
	if err != nil {
		return err
	}

	internalSnap, ok := loaded.(InternalSnapshot)
	if !ok {
		return errors.New("loaded snapshot is not of type InternalSnapshot")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.store = make(map[core.IntentID]*stateEntry)
	for k, v := range internalSnap.Data {
		parsedUUID, err := uuid.Parse(k)
		if err == nil {
			r.store[core.IntentID(parsedUUID)] = v
		}
	}

	// In a real system, we would then replay the ledger from `snapshot.Offset`
	// Since History Recorder is out of scope for the Repository contract directly,
	// the App/Runtime bootstrap phase coordinates Ledger replay *after* calling Recover.

	return nil
}
