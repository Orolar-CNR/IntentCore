package core

import "github.com/google/uuid"

// IntentID uniquely identifies an Intent in the system (UUIDv4).
type IntentID uuid.UUID

// StateVersion represents the monotonically increasing state version.
type StateVersion uint64

// TraceID is used for tracking requests through the system.
type TraceID string

// Authority represents the level of authority invoking a transition.
type Authority string
