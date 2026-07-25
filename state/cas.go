package state

import (
	"github.com/Orolar-CNR/IntentCore/core"
)

// EvaluateCAS compares the expected version with the current version.
// If they do not match, it returns a core.ErrVersionConflict.
func EvaluateCAS(expected, current core.StateVersion) error {
	if expected != current {
		return core.ErrVersionConflict
	}
	return nil
}
