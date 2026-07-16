package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/history"
	"github.com/Orolar-CNR/IntentCore/lifecycle"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/Orolar-CNR/IntentCore/transport"
)

func main() {
	// 1. Initialize In-memory State Repository & History Ledger
	repo := state.NewRepository()
	ledger := history.NewLedger()
	recorder := history.NewRecorder(ledger)

	// 2. Initialize Lifecycle Machine
	machine := lifecycle.NewStateMachine(repo, recorder)

	// 3. Initialize Admission Evaluator
	evaluator := admission.NewPolicyEvaluator()

	// 4. Initialize Transport (Mock for Phase 1)
	// We'll pass an empty list of payloads just to satisfy the DI.
	// Actual E2E testing will happen in the test file.
	trans := transport.NewMockTransport([][]byte{})

	// 5. Initialize Runtime Pipeline
	pipeline := runtime.NewPipeline(trans, evaluator, machine)

	fmt.Println("Starting IntentCore runnable skeleton...")

	if err := pipeline.Start(context.Background()); err != nil {
		log.Fatalf("Pipeline execution failed: %v", err)
	}

	fmt.Println("IntentCore shutdown successfully.")
}
