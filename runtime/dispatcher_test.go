package runtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/google/uuid"
)

// mockLifecycle implements contracts.Lifecycle for testing
type mockLifecycle struct {
	transitions int
	lastRequest contracts.TransitionRequest
	done        chan struct{}
}

func (m *mockLifecycle) Transition(ctx context.Context, req contracts.TransitionRequest) error {
	m.transitions++
	m.lastRequest = req
	if m.done != nil {
		m.done <- struct{}{}
	}
	return nil
}

func TestDispatcher(t *testing.T) {
	ml := &mockLifecycle{
		done: make(chan struct{}, 1),
	}

	// Create dispatcher with small queue
	d := runtime.NewDispatcher(ml, 2)
	d.Start()

	env := contracts.SemanticEnvelope{
		EnvelopeID:    core.IntentID(uuid.New()),
		AgentIdentity: "agent-1",
		OpaquePayload: []byte("payload"),
	}

	// Dispatch should succeed
	err := d.Dispatch(context.Background(), env)
	if err != nil {
		t.Fatalf("Expected dispatch to succeed, got %v", err)
	}

	// Wait for worker to process
	select {
	case <-ml.done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for dispatcher worker")
	}

	if ml.transitions != 1 {
		t.Errorf("Expected 1 transition, got %d", ml.transitions)
	}

	if ml.lastRequest.IntentID != env.EnvelopeID {
		t.Errorf("Expected IntentID %s, got %s", env.EnvelopeID, ml.lastRequest.IntentID)
	}

	// Stop should close down cleanly
	d.Stop()

	// Dispatching after stop should fail
	err = d.Dispatch(context.Background(), env)
	if err == nil {
		t.Error("Expected error dispatching after stop, got nil")
	}
}
