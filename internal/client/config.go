package client

import (
	"time"
)

// Config holds client configuration
type Config struct {
	// Server addresses
	RESTAddr string
	GRPCAddr string

	// Client mode: "grpc" or "rest"
	Mode string

	// Timeout for requests
	Timeout time.Duration

	// Output format: "json" or "table"
	OutputFormat string

	// Enable verbose logging
	Verbose bool
}

// DefaultConfig returns default client configuration
func DefaultConfig() *Config {
	return &Config{
		RESTAddr:     "http://localhost:8080",
		GRPCAddr:     "localhost:9090",
		Mode:         "grpc", // Default to gRPC
		Timeout:      30 * time.Second,
		OutputFormat: "json",
		Verbose:      false,
	}
}
