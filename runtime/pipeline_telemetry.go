package runtime

import (
	"github.com/Orolar-CNR/IntentCore/contracts"
)

// WithTelemetry registers an optional TelemetryRecorder with the pipeline.
func WithTelemetry(recorder contracts.TelemetryRecorder) PipelineOption {
	return func(p *DefaultPipeline) {
		p.telemetry = recorder
	}
}
