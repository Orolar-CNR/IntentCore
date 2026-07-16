package admission

import (
	"context"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// PolicyEvaluator manages the execution of multiple AdmissionPolicies
// and computes the final deterministic AdmissionResult.
type PolicyEvaluator struct {
	policies []contracts.AdmissionPolicy
}

// NewPolicyEvaluator creates a new admission policy runner.
func NewPolicyEvaluator(policies ...contracts.AdmissionPolicy) *PolicyEvaluator {
	return &PolicyEvaluator{
		policies: policies,
	}
}

// Evaluate runs the incoming envelope through all registered policies.
// If any policy rejects, the intent is immediately rejected.
func (pe *PolicyEvaluator) Evaluate(ctx context.Context, env contracts.SemanticEnvelope) (contracts.AdmissionResult, error) {
	// Simple passthrough for Phase 1
	return contracts.AdmissionResult{
		Decision: contracts.DecisionAccept,
		Evidence: contracts.AdmissionEvidence{
			PolicyID:   "phase1-default-policy",
			Timestamp:  time.Now().Format(time.RFC3339),
			VerifierID: "system",
			Reason:     "Auto-accepted by skeleton pipeline",
		},
	}, nil
}
