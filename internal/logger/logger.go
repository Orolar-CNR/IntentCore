package logger

import (
	"log/slog"
	"os"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// SlogLogger is an implementation of contracts.Logger using the standard log/slog package.
type SlogLogger struct {
	l *slog.Logger
}

// New creates a new SlogLogger with the specified log level.
func New(level string) *SlogLogger {
	var programLevel = new(slog.LevelVar) // Info by default

	switch level {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelInfo)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slogger := slog.New(h)

	return &SlogLogger{
		l: slogger,
	}
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

var _ contracts.Logger = (*SlogLogger)(nil)
