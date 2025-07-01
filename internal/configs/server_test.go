package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig_Defaults(t *testing.T) {
	cfg := NewServerConfig()

	assert.Equal(t, "localhost:8080", cfg.Address)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestNewServerConfig_WithAddress(t *testing.T) {
	cfg := NewServerConfig(
		WithServerAddress("127.0.0.1:9000"),
	)

	assert.Equal(t, "127.0.0.1:9000", cfg.Address)
	assert.Equal(t, "info", cfg.LogLevel) // default unchanged
}

func TestNewServerConfig_WithLogLevel(t *testing.T) {
	cfg := NewServerConfig(
		WithServerLogLevel("debug"),
	)

	assert.Equal(t, "localhost:8080", cfg.Address) // default unchanged
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestNewServerConfig_WithMultipleOpts(t *testing.T) {
	cfg := NewServerConfig(
		WithServerAddress("0.0.0.0:1234"),
		WithServerLogLevel("warn"),
	)

	assert.Equal(t, "0.0.0.0:1234", cfg.Address)
	assert.Equal(t, "warn", cfg.LogLevel)
}
