package runtime

import (
	"github.com/Orolar-CNR/IntentCore/contracts"
)

// WithProof registers an optional ProofRecorder with the pipeline.
func WithProof(recorder contracts.ProofRecorder) PipelineOption {
	return func(p *DefaultPipeline) {
		p.proof = recorder
	}
}
