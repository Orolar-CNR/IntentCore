package runtime

import (
	"encoding/json"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// ValidateEnvelope performs basic structural validation on the raw payload.
// For Phase 1, we assume the payload is JSON-encoded SemanticEnvelope.
func ValidateEnvelope(payload []byte) (contracts.SemanticEnvelope, error) {
	var env contracts.SemanticEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return env, core.ErrValidationFailed
	}

	// Basic structural checks could go here (e.g. UUID empty check)
	// if env.EnvelopeID == uuid.Nil { return env, core.ErrValidationFailed }

	return env, nil
}
