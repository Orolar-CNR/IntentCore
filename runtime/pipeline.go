package runtime

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// DefaultPipeline orchestrates the one-way dependency execution model.
type DefaultPipeline struct {
	transport contracts.Transport
	evaluator contracts.AdmissionPolicy
	lifecycle contracts.Lifecycle
}

// NewPipeline creates the root orchestrator for the IntentCore kernel.
func NewPipeline(
	transport contracts.Transport,
	evaluator contracts.AdmissionPolicy,
	lifecycle contracts.Lifecycle,
) *DefaultPipeline {
	return &DefaultPipeline{
		transport: transport,
		evaluator: evaluator,
		lifecycle: lifecycle,
	}
}

// Execute processes a raw payload through the RFC-0000 pipeline:
// Validation -> Normalization -> Admission -> Lifecycle establishes Pending/advances state
func (p *DefaultPipeline) Execute(ctx context.Context, rawPayload []byte) error {
	// 1. Validation
	env, err := ValidateEnvelope(rawPayload)
	if err != nil {
		return err
	}

	// 2. Normalization
	env = NormalizeEnvelope(env)

	// 3. Admission
	result, err := p.evaluator.Evaluate(ctx, env)
	if err != nil {
		return err
	}
	if result.Decision == contracts.DecisionReject {
		return core.ErrAdmissionRejected
	}

	// 4. Lifecycle (establishes Pending / advances state)
	req := contracts.TransitionRequest{
		IntentID:        env.EnvelopeID,
		FromState:       "", // Indicates a new intent
		ToState:         contracts.StatePending,
		ExpectedVersion: 0,
		ActorID:         env.AgentIdentity,
		Metadata: map[string]any{
			"AgentIdentity": env.AgentIdentity,
			"Payload":       env.OpaquePayload,
		},
	}

	if err := p.lifecycle.Transition(ctx, req); err != nil {
		return err
	}

	return nil
}

// Start begins the transport layer, feeding data into Execute.
func (p *DefaultPipeline) Start(ctx context.Context) error {
	return p.transport.Start(ctx, p.Execute)
}
