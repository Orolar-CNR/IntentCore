package state

import (
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"time"
)

// stateEntry represents the internal storage structure for an intent in the Repository.
// It maps the public IntentRecord contract into a format suitable for internal storage.
type stateEntry struct {
	ID        core.IntentID
	Agent     string
	Timestamp time.Time
	Payload   []byte

	State   contracts.IntentState
	Version core.StateVersion
}

func (s *stateEntry) toRecord() *contracts.IntentRecord {
	return &contracts.IntentRecord{
		ID:        s.ID,
		Agent:     s.Agent,
		Timestamp: s.Timestamp,
		Payload:   s.Payload,
		State:     s.State,
		Version:   s.Version,
	}
}
