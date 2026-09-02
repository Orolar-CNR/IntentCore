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

type SnapshotStoreExtended interface {
	contracts.SnapshotStore
	SaveInternal(ctx context.Context, snapshot any) error
	LoadLatest(ctx context.Context) (any, error)

	// LoadInternal returns the full InternalSnapshot for a specific
	// snapshot ID, or (nil, nil) if it does not exist. This lets Recover
	// restore an explicitly requested snapshot instead of always the
	// most recent one.
	LoadInternal(ctx context.Context, id string) (any, error)
}

// Repository implements contracts.StateRepository.
type Repository struct {
	mu            sync.RWMutex
	snapshotStore SnapshotStoreExtended
	data          map[uuid.UUID]*stateEntry
}

func NewRepository(store contracts.SnapshotStore) *Repository {
	if store == nil {
		store = NewInMemorySnapshotStore()
	}
	ext, ok := store.(SnapshotStoreExtended)
	if !ok {
		panic("provided SnapshotStore does not implement SnapshotStoreExtended")
	}
	return &Repository{
		snapshotStore: ext,
		data:          make(map[uuid.UUID]*stateEntry),
	}
}

func (r *Repository) LoadIntent(ctx context.Context, id core.IntentID) (*contracts.IntentRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uid := uuid.UUID(id)
	entry, ok := r.data[uid]
	if !ok {
		return nil, core.ErrNotFound
	}

	return entry.toRecord(), nil
}

func (r *Repository) CompareAndSwap(ctx context.Context, expected core.StateVersion, record contracts.IntentRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate revision invariant: record.Version must strictly equal expected + 1
	if record.Version != expected+1 {
		return core.ErrVersionConflict
	}

	uid := uuid.UUID(record.ID)
	current, exists := r.data[uid]

	if !exists {
		if expected != 0 {
			return core.ErrVersionConflict
		}
	} else {
		if current.Version != expected {
			return core.ErrVersionConflict
		}
	}

	r.data[uid] = &stateEntry{
		ID:        uid,
		Agent:     record.Agent,
		State:     record.State,
		Version:   record.Version,
		Timestamp: record.Timestamp,
		Payload:   record.Payload,
	}

	return nil
}

func (r *Repository) Snapshot(ctx context.Context) (*contracts.Snapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := make(map[string]*stateEntry)
	for k, v := range r.data {
		clone := *v
		data[uuid.UUID(k).String()] = &clone
	}

	checksum, err := computeChecksum(data)
	if err != nil {
		return nil, err
	}

	snapID := uuid.New().String()
	internalSnap := InternalSnapshot{
		Header: contracts.Snapshot{
			ID:          snapID,
			Offset:      uint64(time.Now().UnixNano()), // Simplified offset for phase 2
			IntentCount: uint64(len(data)),
			Checksum:    checksum,
			SnapshotID:  snapID, // Keep for backwards compatibility with tests
		},
		Data: data,
	}

	if err := r.snapshotStore.SaveInternal(ctx, internalSnap); err != nil {
		return nil, err
	}

	headerCopy := internalSnap.Header
	return &headerCopy, nil
}

func (r *Repository) Recover(ctx context.Context, snapshot contracts.Snapshot) error {
	// Re-hydrate state from the requested snapshot. If a specific ID was
	// given, honor it exactly instead of silently substituting whatever
	// happens to be "latest" in the store. Callers that genuinely want
	// the latest snapshot pass a zero-value contracts.Snapshot{}.
	var loaded any
	var err error
	if snapshot.ID != "" {
		loaded, err = r.snapshotStore.LoadInternal(ctx, snapshot.ID)
	} else {
		loaded, err = r.snapshotStore.LoadLatest(ctx)
	}
	if err != nil {
		return err
	}
	if loaded == nil {
		return core.ErrSnapshotNotFound
	}

	internalSnap, ok := loaded.(InternalSnapshot)
	if !ok {
		return errors.New("loaded snapshot is not of type InternalSnapshot")
	}

	// Verify integrity before accepting any of this snapshot's data into
	// the live repository. A snapshot saved before this check existed
	// will have an empty Checksum and is allowed through unverified;
	// any snapshot that does carry a Checksum must match exactly.
	if internalSnap.Header.Checksum != "" {
		checksum, err := computeChecksum(internalSnap.Data)
		if err != nil {
			return err
		}
		if checksum != internalSnap.Header.Checksum {
			return core.ErrSnapshotChecksumMismatch
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[uuid.UUID]*stateEntry)
	for k, v := range internalSnap.Data {
		id, err := uuid.Parse(k)
		if err != nil {
			return err
		}
		clone := *v
		r.data[id] = &clone
	}

	return nil
}
