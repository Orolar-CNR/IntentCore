package state_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/google/uuid"
)

// Recover must honor an explicitly requested snapshot ID rather than
// silently falling back to whatever is "latest" in the store.
func TestRepository_Recover_ByID_NotLatest(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo1 := state.NewRepository(store)

	id1 := core.IntentID(uuid.New())
	record1 := contracts.IntentRecord{
		ID:        id1,
		Agent:     "agent-1",
		State:     contracts.StatePending,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   []byte("first"),
	}
	if err := repo1.CompareAndSwap(ctx, 0, record1); err != nil {
		t.Fatalf("CAS 1 failed: %v", err)
	}

	header1, err := repo1.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Snapshot 1 failed: %v", err)
	}
	if header1.IntentCount != 1 {
		t.Fatalf("expected 1 intent in snapshot 1, got %d", header1.IntentCount)
	}

	id2 := core.IntentID(uuid.New())
	record2 := contracts.IntentRecord{
		ID:        id2,
		Agent:     "agent-2",
		State:     contracts.StatePending,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   []byte("second"),
	}
	if err := repo1.CompareAndSwap(ctx, 0, record2); err != nil {
		t.Fatalf("CAS 2 failed: %v", err)
	}

	header2, err := repo1.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Snapshot 2 failed: %v", err)
	}
	if header2.IntentCount != 2 {
		t.Fatalf("expected 2 intents in snapshot 2, got %d", header2.IntentCount)
	}

	// Recover explicitly to the OLDER snapshot (header1). If Recover
	// silently used LoadLatest instead of honoring the given ID, this
	// would incorrectly restore snapshot 2's two intents instead of
	// snapshot 1's one intent.
	repo2 := state.NewRepository(store)
	if err := repo2.Recover(ctx, *header1); err != nil {
		t.Fatalf("Recover(header1) failed: %v", err)
	}

	if _, err := repo2.LoadIntent(ctx, id1); err != nil {
		t.Fatalf("expected id1 present after recovering snapshot 1: %v", err)
	}
	if _, err := repo2.LoadIntent(ctx, id2); !errors.Is(err, core.ErrNotFound) {
		t.Fatalf("expected id2 absent (snapshot 1 predates it), got err=%v", err)
	}
}

// Recover must reject a snapshot whose stored data no longer matches the
// checksum recorded at save time.
func TestRepository_Recover_ChecksumMismatch(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo1 := state.NewRepository(store)

	id := core.IntentID(uuid.New())
	record := contracts.IntentRecord{
		ID:        id,
		Agent:     "test-agent",
		State:     contracts.StatePending,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   []byte("original"),
	}
	if err := repo1.CompareAndSwap(ctx, 0, record); err != nil {
		t.Fatalf("CAS failed: %v", err)
	}

	header, err := repo1.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}
	if header.Checksum == "" {
		t.Fatalf("expected Snapshot() to populate a non-empty Checksum")
	}

	// Simulate tampering/corruption of the snapshot at rest: mutate the
	// stored entry's payload in place without updating the recorded
	// checksum, exactly like a bit-flip or unauthorized edit would.
	loadedRaw, err := store.LoadInternal(ctx, header.ID)
	if err != nil {
		t.Fatalf("LoadInternal failed: %v", err)
	}
	internal, ok := loadedRaw.(state.InternalSnapshot)
	if !ok {
		t.Fatalf("unexpected type %T from LoadInternal", loadedRaw)
	}
	if len(internal.Data) == 0 {
		t.Fatalf("expected at least one entry in snapshot data")
	}
	for _, entry := range internal.Data {
		entry.Payload = []byte("tampered")
	}

	repo2 := state.NewRepository(store)
	err = repo2.Recover(ctx, *header)
	if !errors.Is(err, core.ErrSnapshotChecksumMismatch) {
		t.Fatalf("expected core.ErrSnapshotChecksumMismatch, got %v", err)
	}
}

// Recover must return core.ErrSnapshotNotFound for an explicitly requested
// snapshot ID that does not exist, instead of silently substituting
// whatever is "latest".
func TestRepository_Recover_NotFound(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo := state.NewRepository(store)

	err := repo.Recover(ctx, contracts.Snapshot{ID: "does-not-exist"})
	if !errors.Is(err, core.ErrSnapshotNotFound) {
		t.Fatalf("expected core.ErrSnapshotNotFound, got %v", err)
	}
}

// A snapshot saved before the Checksum field existed (empty string) must
// still be recoverable without being rejected, for backward compatibility.
func TestRepository_Recover_EmptyChecksum_SkipsVerification(t *testing.T) {
	ctx := context.Background()
	store := state.NewInMemorySnapshotStore()
	repo1 := state.NewRepository(store)

	id := core.IntentID(uuid.New())
	record := contracts.IntentRecord{
		ID:        id,
		Agent:     "legacy-agent",
		State:     contracts.StatePending,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   []byte("legacy"),
	}
	if err := repo1.CompareAndSwap(ctx, 0, record); err != nil {
		t.Fatalf("CAS failed: %v", err)
	}

	header, err := repo1.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	// Simulate a pre-patch snapshot on disk by clearing the checksum that
	// was just computed, then persist that back into the store.
	loadedRaw, err := store.LoadInternal(ctx, header.ID)
	if err != nil {
		t.Fatalf("LoadInternal failed: %v", err)
	}
	internal, ok := loadedRaw.(state.InternalSnapshot)
	if !ok {
		t.Fatalf("unexpected type %T from LoadInternal", loadedRaw)
	}
	internal.Header.Checksum = ""
	if err := store.SaveInternal(ctx, internal); err != nil {
		t.Fatalf("SaveInternal failed: %v", err)
	}
	legacyHeader := internal.Header

	repo2 := state.NewRepository(store)
	if err := repo2.Recover(ctx, legacyHeader); err != nil {
		t.Fatalf("expected legacy snapshot with empty Checksum to recover successfully, got: %v", err)
	}
	if _, err := repo2.LoadIntent(ctx, id); err != nil {
		t.Fatalf("expected id present after recovering legacy snapshot: %v", err)
	}
}
