package contracts

import (
	"context"
)

// Pipeline represents the strictly ordered, one-way dependency execution model.
//
// RFC:
//
//	RFC-0000 Section 3.3 (One-Way Dependency)
//
// Pipeline Flow:
// Transport -> SemanticEnvelope -> Validation -> Normalization -> Admission -> Lifecycle -> Repository -> History
type Pipeline interface {
	// Execute processes a raw byte payload through the entire architectural pipeline.
	Execute(ctx context.Context, rawPayload []byte) error
}
