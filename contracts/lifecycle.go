package contracts

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/core"
)

// Lifecycle defines the authority for coordinating intent state transitions.
//
// RFC:
//
//	RFC-0004
//
// Guarantees:
//   - Deterministic execution of the Transition Matrix.
//   - Atomic state transitions utilizing the StateRepository.
type Lifecycle interface {
	// Transition attempts to move an intent to a target state.
	// Returns core.ErrInvalidTransition if the target state is not allowed from the current state.
	// Returns core.ErrTerminalState if the current state is already terminal.
	Transition(ctx context.Context, intent core.IntentID, targetState IntentState) error
}
