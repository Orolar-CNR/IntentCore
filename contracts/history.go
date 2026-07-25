package contracts

import (
	"context"
	"github.com/Orolar-CNR/IntentCore/core"
)

// HistoryRecorder defines the interface for recording state transitions and events
// into an immutable ledger or evidence store.
//
// Guarantees:
//   - Append-only recording.
//   - Must accurately reflect state transitions that occurred.
type HistoryRecorder interface {
	// RecordTransition records a completed state transition.
	RecordTransition(ctx context.Context, req TransitionRequest, finalState IntentState, newVersion core.StateVersion) error
}
