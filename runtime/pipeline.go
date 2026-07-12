package runtime

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// DefaultPipeline orchestrates the one-way dependency execution model.
type DefaultPipeline struct {
	transport contracts.Transport
	evaluator contracts.AdmissionPolicy // Simplification for skeleton
	lifecycle contracts.Lifecycle
	repo      contracts.StateRepository
}

// NewPipeline creates the root orchestrator for the IntentCore kernel.
func NewPipeline(
	transport contracts.Transport,
	evaluator contracts.AdmissionPolicy,
	lifecycle contracts.Lifecycle,
	repo contracts.StateRepository,
) *DefaultPipeline {
	return &DefaultPipeline{
		transport: transport,
		evaluator: evaluator,
		lifecycle: lifecycle,
		repo:      repo,
	}
}

// Execute processes a raw payload through the RFC-0000 pipeline:
// Validation -> Normalization -> Admission -> Lifecycle -> Repository
func (p *DefaultPipeline) Execute(ctx context.Context, rawPayload []byte) error {
	panic("not implemented: execution pipeline coordination")
}

// Start begins the transport layer, feeding data into Execute.
func (p *DefaultPipeline) Start(ctx context.Context) error {
	return p.transport.Start(ctx, p.Execute)
}
