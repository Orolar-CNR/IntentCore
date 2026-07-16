package contracts

// Logger defines the minimal interface for logging across the application.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
