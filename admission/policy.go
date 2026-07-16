package admission

import (
	"context"
	"errors"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/google/uuid"
)

// DefaultSchemaVersion is the expected canonical schema version based on RFC-0000.
const DefaultSchemaVersion = "1.0.0"

// DeterministicPolicy implements contracts.AdmissionPolicy checking deterministic rules.
type DeterministicPolicy struct{}

// Evaluate assesses a validated SemanticEnvelope against RFC-0002 deterministic rules.
func (dp *DeterministicPolicy) Evaluate(ctx context.Context, env contracts.SemanticEnvelope) (contracts.AdmissionResult, error) {
	// 1. Envelope Version mismatch
	if env.SchemaVersion != "" && env.SchemaVersion != DefaultSchemaVersion {
		return dp.reject(contracts.CodeInvalidVersion, "Unsupported SchemaVersion"), nil
	}

	// 2. Missing IntentID
	if env.EnvelopeID == core.IntentID(uuid.Nil) {
		return dp.reject(contracts.CodeMissingIntentID, "IntentID is required"), nil
	}

	// 3. Missing AgentIdentity
	if env.AgentIdentity == "" {
		return dp.reject(contracts.CodeMissingIdentity, "AgentIdentity is required"), nil
	}

	// 4. Missing Signature (Simulating empty slice check for now)
	if len(env.Signatures) == 0 {
		return dp.reject(contracts.CodeMissingSignature, "Signature is required"), nil
	}

	// 5. Missing Payload
	if len(env.OpaquePayload) == 0 {
		return dp.reject(contracts.CodeMissingPayload, "Payload is required"), nil
	}

	// 6. Invalid Timestamp
	if env.EventTimestamp.IsZero() {
		return dp.reject(contracts.CodeInvalidTimestamp, "Valid EventTimestamp is required"), nil
	}

	return contracts.AdmissionResult{
		Decision: contracts.DecisionAccept,
		Evidence: contracts.AdmissionEvidence{
			PolicyID:   "deterministic-policy",
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			VerifierID: "system",
			Reason:     "Envelope satisfies all deterministic rules",
		},
	}, nil
}

func (dp *DeterministicPolicy) reject(code contracts.RejectionCode, reason string) contracts.AdmissionResult {
	return contracts.AdmissionResult{
		Decision: contracts.DecisionReject,
		Evidence: contracts.AdmissionEvidence{
			PolicyID:   "deterministic-policy",
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			VerifierID: "system",
			Reason:     reason,
			Code:       code,
		},
		Error: errors.New(reason),
	}
}

// PolicyEvaluator manages the execution of multiple AdmissionPolicies
// and computes the final deterministic AdmissionResult.
type PolicyEvaluator struct {
	policies []contracts.AdmissionPolicy
}

// NewPolicyEvaluator creates a new admission policy runner.
// It injects the DeterministicPolicy by default.
func NewPolicyEvaluator(policies ...contracts.AdmissionPolicy) *PolicyEvaluator {
	if len(policies) == 0 {
		policies = []contracts.AdmissionPolicy{&DeterministicPolicy{}}
	}
	return &PolicyEvaluator{
		policies: policies,
	}
}

// Evaluate runs the incoming envelope through all registered policies.
// If any policy rejects, the intent is immediately rejected.
func (pe *PolicyEvaluator) Evaluate(ctx context.Context, env contracts.SemanticEnvelope) (contracts.AdmissionResult, error) {
	for _, policy := range pe.policies {
		res, err := policy.Evaluate(ctx, env)
		if err != nil {
			return res, err
		}
		if res.Decision == contracts.DecisionReject {
			return res, nil
		}
	}

	return contracts.AdmissionResult{
		Decision: contracts.DecisionAccept,
		Evidence: contracts.AdmissionEvidence{
			PolicyID:   "composite-evaluator",
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			VerifierID: "system",
			Reason:     "Auto-accepted by all policies",
		},
	}, nil
}
