package lifecycle

import (
	"github.com/Orolar-CNR/IntentCore/contracts"
)

// CheckAuthority evaluates if the given transition request has the necessary
// authority to execute the transition. This is a placeholder for more
// complex RBAC or rule-based authority checks.
func CheckAuthority(req contracts.TransitionRequest) bool {
	// For Phase 1, we simply return true. Real implementations would verify
	// the requested authority against the actor ID and target state.
	return true
}
