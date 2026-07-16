package contracts

import (
	"github.com/Orolar-CNR/IntentCore/core"
	"time"
)

// IntentState represents the current state of an Intent in the lifecycle.
//
// RFC:
//
//	RFC-0004 Section 3 (Allowed Transition Matrix)
type IntentState string

const (
	StatePending    IntentState = "Pending"
	StateValidated  IntentState = "Validated"
	StateAdmitted   IntentState = "Admitted"
	StateScheduled  IntentState = "Scheduled"
	StateExecuting  IntentState = "Executing"
	StateCompleted  IntentState = "Completed"
	StateFailed     IntentState = "Failed"
	StateRolledBack IntentState = "RolledBack"
)

// IntentRecord represents the canonical form of an intent as stored in the Repository.
// It bundles the SemanticEnvelope data with its current Lifecycle state and version.
type IntentRecord struct {
	ID        core.IntentID
	Agent     string
	Timestamp time.Time
	Payload   []byte // JSON encoded opaque payload

	State   IntentState
	Version core.StateVersion
}
