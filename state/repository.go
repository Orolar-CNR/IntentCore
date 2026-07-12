package state

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// Repository implements contracts.StateRepository.
// It maintains the single source of truth for the system, enforcing CAS semantics.
type Repository struct {
	// Implementation details (e.g. database connection, in-memory store) go here.
}

// NewRepository initializes a new State Repository.
func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) LoadIntent(ctx context.Context, id core.IntentID) (*contracts.IntentRecord, error) {
	panic("not implemented: see RFC-0003")
}

func (r *Repository) CompareAndSwap(ctx context.Context, expected core.Version, next contracts.IntentRecord) error {
	panic("not implemented: see RFC-0003")
}

func (r *Repository) Snapshot(ctx context.Context) (*contracts.Snapshot, error) {
	panic("not implemented: see RFC-0003")
}

func (r *Repository) Recover(ctx context.Context, snapshot contracts.Snapshot) error {
	panic("not implemented: see RFC-0003")
}
