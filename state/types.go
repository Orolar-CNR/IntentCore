package state

import (
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/google/uuid"
)

// stateEntry represents the internal storage structure for an intent in the Repository.
// It maps the public IntentRecord contract into a format suitable for internal storage.
type stateEntry struct {
	ID        uuid.UUID
	Agent     string
	State     contracts.IntentState
	Version   core.StateVersion
	Timestamp time.Time
	Payload   []byte
}

func (s *stateEntry) toRecord() *contracts.IntentRecord {
	return &contracts.IntentRecord{
		ID:        core.IntentID(s.ID),
		Agent:     s.Agent,
		State:     s.State,
		Version:   s.Version,
		Timestamp: s.Timestamp,
		Payload:   s.Payload,
	}
}
