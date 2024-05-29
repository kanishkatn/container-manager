package config

import (
	"testing"
)

func TestConfig_Validate_WithValidConfig(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		JRPCPort:      8080,
		P2PPort:       4001,
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
		JRPCPort:      8080,
		P2PPort:       4001,
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
		JRPCPort:      8080,
		P2PPort:       4001,
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
		JRPCPort:      8080,
		P2PPort:       4001,
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
		JRPCPort:      0,
		P2PPort:       4001,
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
		JRPCPort:      8080,
		P2PPort:       4001,
		LogLevel:      "",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestConfig_ValidateWithEmptyP2PPort(t *testing.T) {
	c := &Config{
		QueueSize:     100,
		WorkerCount:   10,
		ListenAddress: "0.0.0.0",
		JRPCPort:      8080,
		P2PPort:       0,
		LogLevel:      "info",
	}
	err := c.ValidateBasic()
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
