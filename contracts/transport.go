package contracts

import (
	"context"
)

// EnvelopeHandler is a callback provided by the runtime to process envelopes
// successfully received and framed by the Transport.
type EnvelopeHandler func(ctx context.Context, data []byte) error

// Transport defines the boundary interface for delivering data into IntentCore.
//
// RFC:
//
//	RFC-0000 Section 3.5 (Transport Independence)
//	RFC-0001
//
// Guarantees:
//   - Must only handle framing, deserialization, and delivery.
//   - Must not evaluate intent semantics or govern lifecycle.
type Transport interface {
	// Start begins listening for incoming frames on the transport mechanism.
	// It calls handler for each successfully framed byte payload.
	Start(ctx context.Context, handler EnvelopeHandler) error

	// Stop gracefully terminates the transport listener.
	Stop(ctx context.Context) error
}
