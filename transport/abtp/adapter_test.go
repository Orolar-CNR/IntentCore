package abtp_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/transport/abtp"
)

func TestAdapter(t *testing.T) {
	// Pick a random port for testing
	adapter := abtp.NewAdapter(0) // 0 lets OS pick random port, but wait, our adapter hardcodes DefaultPort if <= 0

	// Let's create an adapter on an explicit high port for testing
	adapter = abtp.NewAdapter(10055)

	received := make(chan []byte, 1)

	// Simple handler
	handler := func(ctx context.Context, payload []byte) error {
		received <- payload
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := adapter.Start(ctx, handler)
	if err != nil {
		t.Fatalf("Failed to start adapter: %v", err)
	}
	defer adapter.Stop(context.Background())

	// Send a test packet
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10055")
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	payload := []byte("test-payload")
	_, err = conn.Write(payload)
	if err != nil {
		t.Fatalf("Failed to write UDP: %v", err)
	}

	// Wait for receipt
	select {
	case p := <-received:
		if string(p) != string(payload) {
			t.Errorf("Expected payload %s, got %s", payload, p)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for payload")
	}
}
