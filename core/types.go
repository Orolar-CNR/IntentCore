package core

import "github.com/google/uuid"

// IntentID uniquely identifies an Intent in the system (UUIDv4).
type IntentID uuid.UUID

// Version represents the monotonically increasing state version.
type Version uint64
