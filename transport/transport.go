package transport

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// ABTPTransport wraps the underlying ABTP/eBPF implementation to conform
// to the architectural contracts.Transport interface.
//
// RFC:
//
//	RFC-0000 Section 3.5 (Transport Independence)
type ABTPTransport struct {
	// internal state for ABTP socket/link
}

// NewABTPTransport initializes the ABTP transport mechanism.
func NewABTPTransport() *ABTPTransport {
	return &ABTPTransport{}
}

// Start begins the eBPF/XDP listener and passes validated frames to the handler.
func (t *ABTPTransport) Start(ctx context.Context, handler contracts.EnvelopeHandler) error {
	panic("not implemented: ABTP internal listener hook")
}

// Stop cleanly detaches the XDP program and closes sockets.
func (t *ABTPTransport) Stop(ctx context.Context) error {
	panic("not implemented")
}

// Ensure ABTPTransport implements the contract
var _ contracts.Transport = (*ABTPTransport)(nil)
