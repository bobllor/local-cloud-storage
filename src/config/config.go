package config

import (
	"github.com/bobllor/gologger"
)

// Config is used to provide utilities for structs.
type Config struct {
	// Log is the logging struct.
	Log *gologger.Logger
}

// NewConfig creates a new Config.
func NewConfig(logger *gologger.Logger) *Config {
	return &Config{
		Log: logger,
	}
}
