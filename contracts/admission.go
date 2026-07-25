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

// RejectionCode defines standardized deterministic error codes for rejection.
type RejectionCode string

const (
	CodeInvalidVersion         RejectionCode = "ERR_INVALID_VERSION"
	CodeMissingSignature       RejectionCode = "ERR_MISSING_SIGNATURE"
	CodeMissingIntentID        RejectionCode = "ERR_MISSING_INTENT_ID"
	CodeMissingIdentity        RejectionCode = "ERR_MISSING_IDENTITY"
	CodeMissingPayload         RejectionCode = "ERR_MISSING_PAYLOAD"
	CodeInvalidTimestamp       RejectionCode = "ERR_INVALID_TIMESTAMP"
	CodeGeneralPolicyViolation RejectionCode = "ERR_POLICY_VIOLATION"
)

// AdmissionEvidence provides an auditable trail of an admission decision.
type AdmissionEvidence struct {
	PolicyID   string
	Timestamp  string
	VerifierID string
	Reason     string
	Code       RejectionCode
}

// AdmissionResult represents the outcome of a policy evaluation along with audit evidence.
type AdmissionResult struct {
	Decision AdmissionDecision
	Evidence AdmissionEvidence
	Error    error
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
