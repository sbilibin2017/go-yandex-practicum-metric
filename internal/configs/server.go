package configs

// ServerConfig holds configuration for the server
type ServerConfig struct {
	Address  string `json:"address"`
	LogLevel string `json:"log_level"`
}

// ServerOpt is a functional option for configuring ServerConfig
type ServerOpt func(*ServerConfig)

// WithAddress sets the server address
func WithServerAddress(addr string) ServerOpt {
	return func(cfg *ServerConfig) {
		cfg.Address = addr
	}
}

// WithLogLevel sets the server log level
func WithServerLogLevel(level string) ServerOpt {
	return func(cfg *ServerConfig) {
		cfg.LogLevel = level
	}
}

// NewServerConfig creates a ServerConfig with optional functional parameters
func NewServerConfig(opts ...ServerOpt) *ServerConfig {
	cfg := &ServerConfig{
		Address:  "localhost:8080",
		LogLevel: "info",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
