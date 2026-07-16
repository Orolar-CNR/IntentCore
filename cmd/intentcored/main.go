package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Orolar-CNR/IntentCore/internal/bootstrap"
)

func main() {
	// Setup context that listens for termination signals
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// Bootstrap application dependencies
	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run application
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application execution failed: %v", err)
	}
}
