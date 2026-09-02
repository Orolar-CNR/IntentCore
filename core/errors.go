package core

import "errors"

var (
	ErrVersionConflict   = errors.New("repository: CAS version conflict")
	ErrNotFound          = errors.New("repository: intent not found")
	ErrTerminalState     = errors.New("lifecycle: intent is in a terminal state")
	ErrInvalidTransition = errors.New("lifecycle: invalid state transition")
	ErrAdmissionRejected = errors.New("admission: intent rejected by policy")
	ErrValidationFailed  = errors.New("validation: structural validation failed")

	// ErrSnapshotChecksumMismatch indicates that a snapshot's recomputed
	// checksum does not match the checksum recorded at save time, meaning
	// the snapshot data was altered or corrupted after it was written.
	ErrSnapshotChecksumMismatch = errors.New("repository: snapshot checksum mismatch")

	// ErrSnapshotNotFound indicates that Recover was asked to restore a
	// specific snapshot ID that does not exist in the SnapshotStore.
	ErrSnapshotNotFound = errors.New("repository: snapshot not found")
)
