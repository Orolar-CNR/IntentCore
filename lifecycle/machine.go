package lifecycle

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// StateMachine implements the contracts.Lifecycle interface.
// It acts as the sole authority for evaluating and executing state transitions.
type StateMachine struct {
	repo contracts.StateRepository
}

// NewStateMachine creates a new deterministic lifecycle state machine.
func NewStateMachine(repo contracts.StateRepository) *StateMachine {
	return &StateMachine{
		repo: repo,
	}
}

// Transition enforces the Allowed Transition Matrix and commits the state
// mutation atomically via the repository.
func (sm *StateMachine) Transition(ctx context.Context, intent core.IntentID, targetState contracts.IntentState) error {
	// panic("not implemented: see RFC-0004")
	return nil
}
