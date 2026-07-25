package core

import "errors"

var (
	ErrVersionConflict   = errors.New("repository: CAS version conflict")
	ErrNotFound          = errors.New("repository: intent not found")
	ErrTerminalState     = errors.New("lifecycle: intent is in a terminal state")
	ErrInvalidTransition = errors.New("lifecycle: invalid state transition")
	ErrAdmissionRejected = errors.New("admission: intent rejected by policy")
	ErrValidationFailed  = errors.New("validation: structural validation failed")
)
