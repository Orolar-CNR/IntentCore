package contracts

import (
	"context"
)

// AdmissionDecision represents the terminal outcome of an admission evaluation.
type AdmissionDecision string

const (
	DecisionAccept AdmissionDecision = "Accept"
	DecisionReject AdmissionDecision = "Reject"
)

// AdmissionEvidence provides an auditable trail of an admission decision.
type AdmissionEvidence struct {
	PolicyID   string
	Timestamp  string
	VerifierID string
	Reason     string
}

// AdmissionResult represents the outcome of a policy evaluation along with audit evidence.
type AdmissionResult struct {
	Decision AdmissionDecision
	Evidence AdmissionEvidence
}

// AdmissionPolicy defines a single, decidable governance rule.
//
// RFC:
//
//	RFC-0002
//
// Guarantees:
//   - Must be deterministic.
//   - Must not mutate state.
//   - Must execute in bounded time (no unbounded recursion).
type AdmissionPolicy interface {
	// Evaluate assesses a validated SemanticEnvelope against this policy.
	Evaluate(ctx context.Context, env SemanticEnvelope) (AdmissionResult, error)
}
