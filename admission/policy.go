package admission

import (
	"context"

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
	panic("not implemented: see RFC-0002")
}
