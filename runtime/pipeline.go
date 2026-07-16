package runtime

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
)

// DefaultPipeline orchestrates the one-way dependency execution model.
type DefaultPipeline struct {
	transport  contracts.Transport
	evaluator  contracts.AdmissionPolicy
	lifecycle  contracts.Lifecycle
	dispatcher *Dispatcher

	// Pipeline Configuration Options
	timeout time.Duration
	retries int

	// Optional Observers
	telemetry contracts.TelemetryRecorder
	proof     contracts.ProofRecorder
}

// PipelineOption defines a functional option for configuring the Pipeline.
type PipelineOption func(*DefaultPipeline)

// WithTimeout sets the maximum duration for executing a single intent up to dispatch.
func WithTimeout(timeout time.Duration) PipelineOption {
	return func(p *DefaultPipeline) {
		p.timeout = timeout
	}
}

// WithRetries sets the number of retries for admission and dispatch (not for side-effect lifecycle execution).
func WithRetries(retries int) PipelineOption {
	return func(p *DefaultPipeline) {
		p.retries = retries
	}
}

// NewPipeline creates the root orchestrator for the IntentCore kernel.
func NewPipeline(
	transport contracts.Transport,
	evaluator contracts.AdmissionPolicy,
	lifecycle contracts.Lifecycle,
	opts ...PipelineOption,
) *DefaultPipeline {
	p := &DefaultPipeline{
		transport: transport,
		evaluator: evaluator,
		lifecycle: lifecycle,
		// Buffer of 100 for basic backpressure
		dispatcher: NewDispatcher(lifecycle, 100),

		// Defaults
		timeout: 5 * time.Second,
		retries: 3,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Execute processes a raw payload through the RFC-0000 pipeline:
// Validation -> Normalization -> Admission -> Dispatcher (Queue -> Worker -> Lifecycle)
// This implements a robust runtime with timeouts, cancellation, and retries.
func (p *DefaultPipeline) Execute(ctx context.Context, rawPayload []byte) error {
	// Respect incoming context, or apply configured pipeline timeout
	execCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// 1. Validation
	env, err := ValidateEnvelope(rawPayload)
	if err != nil {
		return err
	}

	// 2. Normalization
	env = NormalizeEnvelope(env)

	// Retry loop for transient issues during admission and dispatch
	var lastErr error
	for attempt := 0; attempt <= p.retries; attempt++ {
		// Check for context cancellation
		if err := execCtx.Err(); err != nil {
			return err
		}

		// 3. Admission
		result, err := p.evaluator.Evaluate(execCtx, env)
		if err != nil {
			lastErr = err
			// Brief backoff could be added here
			continue
		}

		if result.Decision == contracts.DecisionReject {
			// Deterministic policy rejection is terminal, no retry
			if p.telemetry != nil {
				p.telemetry.Record(execCtx, contracts.TelemetryEvent{
					IntentID:  uuid.UUID(env.EnvelopeID).String(),
					EventType: "Admission_Rejected",
					Timestamp: time.Now(),
					Metadata:  map[string]any{"reason": result.Evidence.Reason},
				})
			}
			return core.ErrAdmissionRejected
		}

		if p.proof != nil {
			p.proof.RecordProof(execCtx, contracts.Proof{
				IntentID: uuid.UUID(env.EnvelopeID).String(),
				ProofID:  "proof-" + uuid.UUID(env.EnvelopeID).String(),
			})
		}

		if p.telemetry != nil {
			p.telemetry.Record(execCtx, contracts.TelemetryEvent{
				IntentID:  uuid.UUID(env.EnvelopeID).String(),
				EventType: "Admission_Accepted",
				Timestamp: time.Now(),
			})
		}

		// 4. Dispatch (establishes Pending / advances state asynchronously via worker)
		if err := p.dispatcher.Dispatch(execCtx, env); err != nil {
			lastErr = err
			continue
		}

		// Success
		return nil
	}

	return lastErr
}

// Start begins the transport layer and dispatcher background workers.
func (p *DefaultPipeline) Start(ctx context.Context) error {
	p.dispatcher.Start()
	return p.transport.Start(ctx, p.Execute)
}

// Stop cleanly shuts down the pipeline and its background workers.
func (p *DefaultPipeline) Stop() {
	// If transport had a Stop(), we'd call it here. (Wait for Phase 2.8 for full App shutdown)
	p.dispatcher.Stop()
}
