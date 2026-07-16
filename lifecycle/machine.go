package lifecycle

import (
	"context"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// StateMachine implements the contracts.Lifecycle interface.
// It acts as the sole authority for evaluating and executing state transitions.
type StateMachine struct {
	repo    contracts.StateRepository
	history contracts.HistoryRecorder
}

// NewStateMachine creates a new deterministic lifecycle state machine.
func NewStateMachine(repo contracts.StateRepository, history contracts.HistoryRecorder) *StateMachine {
	return &StateMachine{
		repo:    repo,
		history: history,
	}
}

// Transition enforces the Allowed Transition Matrix and commits the state
// mutation atomically via the repository.
func (sm *StateMachine) Transition(ctx context.Context, req contracts.TransitionRequest) error {
	// 1. Validate the transition matrix
	if !IsAllowed(req.FromState, req.ToState) {
		if req.FromState == contracts.StateCompleted || req.FromState == contracts.StateRolledBack {
			return core.ErrTerminalState
		}
		return core.ErrInvalidTransition
	}

	// 2. Validate Authority
	if !CheckAuthority(req) {
		return core.ErrInvalidTransition // Or a specific ErrUnauthorized error
	}

	// Fetch current state from repo if needed to preserve fields
	var currentRecord *contracts.IntentRecord
	if req.ExpectedVersion > 0 {
		var err error
		currentRecord, err = sm.repo.LoadIntent(ctx, req.IntentID)
		if err != nil {
			return err
		}
	} else {
		// New intent
		currentRecord = &contracts.IntentRecord{
			ID:        req.IntentID,
			Timestamp: time.Now(),
		}
	}

	// 3. Prepare the new state record
	nextVersion := req.ExpectedVersion + 1
	nextRecord := contracts.IntentRecord{
		ID:        currentRecord.ID,
		Agent:     currentRecord.Agent,
		Timestamp: currentRecord.Timestamp,
		Payload:   currentRecord.Payload,
		State:     req.ToState,
		Version:   nextVersion,
	}

	// If metadata contains payload and agent, we can populate it on first transition (e.g., to Pending)
	if req.ExpectedVersion == 0 {
		if agent, ok := req.Metadata["AgentIdentity"].(string); ok {
			nextRecord.Agent = agent
		}
		if payload, ok := req.Metadata["Payload"].([]byte); ok {
			nextRecord.Payload = payload
		}
	}

	// 4. Commit via CAS
	if err := sm.repo.CompareAndSwap(ctx, req.ExpectedVersion, nextRecord); err != nil {
		return err
	}

	// 5. Emit History
	if sm.history != nil {
		if err := sm.history.RecordTransition(ctx, req, req.ToState, nextVersion); err != nil {
			// In a robust system, history failure might require retry or specific handling,
			// but for this phase we might just return the error.
			return err
		}
	}

	return nil
}
