package contracts

import "context"

// RetentionDecision describes what to do with a snapshot after creation.
type RetentionDecision uint8

const (
	RetentionKeep RetentionDecision = iota
	RetentionArchive
	RetentionDelete
)

// RetentionPolicy decides lifecycle actions for snapshots.
type RetentionPolicy interface {
	Decide(ctx context.Context, snapshot *Snapshot) (RetentionDecision, error)
}
