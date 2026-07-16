package history

import (
	"context"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// Recorder implements contracts.HistoryRecorder, writing to an append-only ledger.
type Recorder struct {
	ledger *Ledger
}

// NewRecorder creates a new History Recorder backed by the given ledger.
func NewRecorder(ledger *Ledger) *Recorder {
	return &Recorder{
		ledger: ledger,
	}
}

// RecordTransition translates a lifecycle transition into an immutable evidence record.
func (r *Recorder) RecordTransition(ctx context.Context, req contracts.TransitionRequest, finalState contracts.IntentState, newVersion core.StateVersion) error {
	record := EvidenceRecord{
		IntentID:        req.IntentID,
		FromState:       req.FromState,
		ToState:         finalState,
		RecordedVersion: newVersion,
		ActorID:         req.ActorID,
		TraceID:         req.TraceID,
		Timestamp:       time.Now(),
	}

	r.ledger.Append(record)
	return nil
}
