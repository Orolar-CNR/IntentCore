package app

import (
	"context"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/internal/config"
)

// Pipeline represents the core execution loop.
type Pipeline interface {
	Start(ctx context.Context) error
	Stop()
}

// App is the runtime wrapper for IntentCore.
type App struct {
	pipeline Pipeline
	logger   contracts.Logger
	cfg      *config.Config
}

// New creates a new App instance.
func New(pipeline Pipeline, logger contracts.Logger, cfg *config.Config) *App {
	return &App{
		pipeline: pipeline,
		logger:   logger,
		cfg:      cfg,
	}
}

// Run starts the application and blocks until the context is canceled.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting IntentCore pipeline...")

	// Start the pipeline in the background
	errCh := make(chan error, 1)
	go func() {
		if err := a.pipeline.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	// Wait for context cancellation or pipeline failure
	select {
	case <-ctx.Done():
		a.logger.Info("Shutdown signal received, gracefully stopping IntentCore...")
		a.pipeline.Stop()
		return nil
	case err := <-errCh:
		a.logger.Error("Pipeline execution failed", "error", err)
		return err
	}
}
