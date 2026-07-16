package transport

import (
	"context"
	"github.com/Orolar-CNR/IntentCore/contracts"
)

// MockTransport is a simple adapter to feed byte payloads into the system for testing.
type MockTransport struct {
	payloads [][]byte
}

func NewMockTransport(payloads [][]byte) *MockTransport {
	return &MockTransport{payloads: payloads}
}

func (t *MockTransport) Start(ctx context.Context, handler contracts.EnvelopeHandler) error {
	for _, payload := range t.payloads {
		if err := handler(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}

func (t *MockTransport) Stop(ctx context.Context) error {
	return nil
}

var _ contracts.Transport = (*MockTransport)(nil)
