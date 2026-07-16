package history

import (
	"sync"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// EvidenceRecord represents an immutable entry in the history ledger.
type EvidenceRecord struct {
	IntentID        core.IntentID
	FromState       contracts.IntentState
	ToState         contracts.IntentState
	RecordedVersion core.StateVersion
	ActorID         string
	TraceID         core.TraceID
	Timestamp       time.Time
}

// Ledger is an append-only in-memory storage of evidence.
type Ledger struct {
	mu      sync.RWMutex
	records []EvidenceRecord
}

// NewLedger creates a new append-only history ledger.
func NewLedger() *Ledger {
	return &Ledger{
		records: make([]EvidenceRecord, 0),
	}
}

// Append adds a new evidence record to the ledger.
func (l *Ledger) Append(record EvidenceRecord) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
}

// GetRecords returns a copy of all records for an intent.
func (l *Ledger) GetRecords(intentID core.IntentID) []EvidenceRecord {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var result []EvidenceRecord
	for _, rec := range l.records {
		if rec.IntentID == intentID {
			result = append(result, rec)
		}
	}
	return result
}
