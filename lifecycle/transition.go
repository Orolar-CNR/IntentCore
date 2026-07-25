package lifecycle

import "github.com/Orolar-CNR/IntentCore/contracts"

// IsAllowed checks if a transition from current to target state is permitted
// by the RFC-0004 Allowed Transition Matrix.
func IsAllowed(current, target contracts.IntentState) bool {
	switch current {
	case contracts.StatePending:
		return target == contracts.StateValidated
	case contracts.StateValidated:
		return target == contracts.StateAdmitted
	case contracts.StateAdmitted:
		return target == contracts.StateScheduled
	case contracts.StateScheduled:
		return target == contracts.StateExecuting
	case contracts.StateExecuting:
		return target == contracts.StateCompleted || target == contracts.StateFailed
	case contracts.StateFailed:
		return target == contracts.StateRolledBack
	case contracts.StateCompleted, contracts.StateRolledBack:
		return false // Terminal states
	default:
		// When no current state is provided (e.g. initial transition)
		if current == "" {
			return target == contracts.StatePending
		}
		return false
	}
}
