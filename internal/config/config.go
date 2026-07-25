package config

// Config holds the application configuration.
type Config struct {
	LogLevel      string
	TransportType string
}

// Default returns a default configuration for Phase 1.
func Default() *Config {
	return &Config{
		LogLevel:      "info",
		TransportType: "mock",
	}
}
