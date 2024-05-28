package config

import (
	"testing"
)

func TestConfig_Validate_WithValidConfig(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		Port:          "8080",
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestConfig_Validate_WithNegativeQueueSize(t *testing.T) {
	c := &Config{
		QueueSize:     -1,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		Port:          "8080",
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestConfig_Validate_WithZeroWorkerCount(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   0,
		ListenAddress: "0.0.0.0",
		Port:          "8080",
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestConfig_Validate_WithEmptyListenAddress(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "",
		Port:          "8080",
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestConfig_Validate_WithEmptyPort(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		Port:          "",
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestConfig_Validate_WithEmptyLogLevel(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		Port:          "8080",
		LogLevel:      "",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
