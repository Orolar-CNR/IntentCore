package contracts

import (
	"time"

	"github.com/Orolar-CNR/IntentCore/core"
)

// SemanticEnvelope represents the canonical wire contract.
//
// RFC:
//
//	RFC-0001 Section 1 (Canonical Schema)
//
// Guarantees:
//   - Must contain a valid UUIDv4 EnvelopeID
//   - Must contain an AgentIdentity string
//   - Must contain an ISO8601 EventTimestamp
//   - Must contain a JSON encoded OpaquePayload
type SemanticEnvelope struct {
	EnvelopeID     core.IntentID
	AgentIdentity  string
	EventTimestamp time.Time
	TelemetryClass string
	OpaquePayload  []byte
}
