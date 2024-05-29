package config

import "fmt"

// Config is the configuration for the container manager
type Config struct {
	// The size of the job queue
	QueueSize int
	// The number of workers to run
	WorkerCount int
	// The address to listen on
	ListenAddress string
	// The JRPCPort to listen on
	JRPCPort int
	// The P2PPort to listen on
	P2PPort int
	// The log level
	LogLevel string
}

// ValidateBasic a basic validation of the config
func (c *Config) ValidateBasic() error {
	if c.QueueSize <= 0 {
		return fmt.Errorf("queue size must be greater than 0")
	}
	if c.WorkerCount <= 0 {
		return fmt.Errorf("worker count must be greater than 0")
	}
	if c.ListenAddress == "" {
		return fmt.Errorf("listen address is required")
	}
	if c.JRPCPort <= 0 {
		return fmt.Errorf("jrpc-port is required")
	}
	if c.P2PPort <= 0 {
		return fmt.Errorf("p2p-port is required")
	}
	if c.LogLevel == "" {
		return fmt.Errorf("log level is required")
	}
	return nil
}

// DefaultConfig returns the default config
func DefaultConfig() *Config {
	return &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		JRPCPort:      8080,
		P2PPort:       4001,
		LogLevel:      "info",
	}
}
