package contracts

import (
	"time"

	"github.com/Orolar-CNR/IntentCore/core"
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
	Version uint64
}

// Intent represents the high-level semantic intent interface.
// This interface defines how the system interacts with an intent internally.
type Intent interface {
	// ID returns the unique identifier of the intent.
	ID() core.IntentID
	// Agent returns the identity of the sender.
	Agent() string
	// Timestamp returns the logical event time.
	Timestamp() time.Time
	// State returns the current lifecycle state of the intent.
	State() IntentState
	// Version returns the current repository version for CAS operations.
	Version() uint64
}
