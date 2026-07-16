package contracts

import (
	"context"
	"time"
)

// ArchivedSnapshot represents an archived snapshot entry.
type ArchivedSnapshot struct {
	ID         string
	SnapshotID string
	CreatedAt  time.Time
	Retention  string
	Location   string
	Payload    []byte
}

// ArchiveStore defines long-term storage for cold snapshots.
type ArchiveStore interface {
	Archive(ctx context.Context, entry *ArchivedSnapshot) error
	Retrieve(ctx context.Context, id string) (*ArchivedSnapshot, error)
	List(ctx context.Context) ([]ArchivedSnapshot, error)
	Purge(ctx context.Context, id string) error
}
