package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// DispatchRequest contains the evaluated envelope ready for execution.
type DispatchRequest struct {
	Envelope contracts.SemanticEnvelope
	Ctx      context.Context
}

// Dispatcher manages the execution queue and routes intent requests to workers.
//
// RFC:
//
//	Phase 2.1 (Runtime Dispatch Core)
//
// Guarantees:
//   - Route and queue evaluated intents.
//   - Isolated execution through workers.
//   - Clean shutdown support.
type Dispatcher struct {
	lifecycle contracts.Lifecycle
	queue     chan DispatchRequest
	wg        sync.WaitGroup
	quit      chan struct{}
	stopped   int32
}

// NewDispatcher creates a new Dispatcher with the given lifecycle engine.
// It initializes an internal channel queue of the specified capacity.
func NewDispatcher(lifecycle contracts.Lifecycle, queueCapacity int) *Dispatcher {
	return &Dispatcher{
		lifecycle: lifecycle,
		queue:     make(chan DispatchRequest, queueCapacity),
		quit:      make(chan struct{}),
	}
}

// Start begins the background worker goroutine(s).
func (d *Dispatcher) Start() {
	d.wg.Add(1)
	go d.worker()
}

// Stop gracefully shuts down the dispatcher, waiting for ongoing work to finish.
func (d *Dispatcher) Stop() {
	if atomic.CompareAndSwapInt32(&d.stopped, 0, 1) {
		close(d.quit)
		d.wg.Wait()
		close(d.queue) // Safe to close queue now that workers are done
	}
}

// Dispatch enqueues an envelope for processing.
// Returns an error if the queue is full or the dispatcher is shutting down.
func (d *Dispatcher) Dispatch(ctx context.Context, env contracts.SemanticEnvelope) error {
	if atomic.LoadInt32(&d.stopped) == 1 {
		return errors.New("dispatcher is shutting down")
	}

	select {
	case <-d.quit:
		return errors.New("dispatcher is shutting down")
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Attempt to enqueue without blocking indefinitely
		select {
		case d.queue <- DispatchRequest{Envelope: env, Ctx: ctx}:
			return nil
		case <-d.quit:
			return errors.New("dispatcher is shutting down")
		default:
			return errors.New("execution queue is full")
		}
	}
}

// worker is the simple consumer loop executing transitions.
func (d *Dispatcher) worker() {
	defer d.wg.Done()

	for {
		select {
		case <-d.quit:
			return
		case req, ok := <-d.queue:
			if !ok {
				return
			}
			// Prepare the transition request
			tReq := contracts.TransitionRequest{
				IntentID:        req.Envelope.EnvelopeID,
				FromState:       "", // indicates new
				ToState:         contracts.StatePending,
				ExpectedVersion: 0,
				ActorID:         req.Envelope.AgentIdentity,
				Metadata: map[string]any{
					"AgentIdentity": req.Envelope.AgentIdentity,
					"Payload":       req.Envelope.OpaquePayload,
				},
			}

			// Execute transition
			err := d.lifecycle.Transition(req.Ctx, tReq)
			if err != nil {
				// In a real system, we'd log this or send to a rejection/dead-letter queue.
				// For now, we print for local debugging.
				fmt.Printf("Worker error executing intent %s: %v\n", req.Envelope.EnvelopeID, err)
			}
		}
	}
}
