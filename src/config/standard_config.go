package config

import (
	"github.com/bobllor/gologger"
)

type StandardConfig struct {
	Log *gologger.Logger
}

func NewStandardConfig(logger *gologger.Logger) *StandardConfig {
	return &StandardConfig{
		Log: logger,
	}
}
