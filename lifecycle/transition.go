package lifecycle

import "github.com/Orolar-CNR/IntentCore/contracts"

// IsAllowed checks if a transition from current to target state is permitted
// by the RFC-0004 Allowed Transition Matrix.
func IsAllowed(current, target contracts.IntentState) bool {
	// panic("not implemented")
	return false
}
