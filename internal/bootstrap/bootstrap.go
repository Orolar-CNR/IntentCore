package bootstrap

import (
	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/history"
	"github.com/Orolar-CNR/IntentCore/internal/app"
	"github.com/Orolar-CNR/IntentCore/internal/config"
	"github.com/Orolar-CNR/IntentCore/internal/logger"
	"github.com/Orolar-CNR/IntentCore/lifecycle"
	"github.com/Orolar-CNR/IntentCore/runtime"
	"github.com/Orolar-CNR/IntentCore/state"
	"github.com/Orolar-CNR/IntentCore/transport"
)

// New initializes and wires all dependencies for the IntentCore application.
func New() (*app.App, error) {
	// 1. Configuration
	cfg := config.Default()

	// 2. Logger
	log := logger.New(cfg.LogLevel)

	// 3. Initialize In-memory State Repository & History Ledger
	repo := state.NewRepository()
	ledger := history.NewLedger()
	recorder := history.NewRecorder(ledger)

	// 4. Initialize Lifecycle Machine
	machine := lifecycle.NewStateMachine(repo, recorder)

	// 5. Initialize Admission Evaluator
	evaluator := admission.NewPolicyEvaluator()

	// 6. Initialize Transport (Mock for Phase 1)
	// Passing an empty list of payloads just to satisfy the DI.
	trans := transport.NewMockTransport([][]byte{})

	// 7. Initialize Runtime Pipeline
	pipeline := runtime.NewPipeline(trans, evaluator, machine)

	// 8. Create App wrapper
	application := app.New(pipeline, log, cfg)

	return application, nil
}
