package runtime

import "github.com/Orolar-CNR/IntentCore/contracts"

// NormalizeEnvelope applies default values and standardizes fields.
func NormalizeEnvelope(env contracts.SemanticEnvelope) contracts.SemanticEnvelope {
	if env.TelemetryClass == "" {
		env.TelemetryClass = "default"
	}
	return env
}
