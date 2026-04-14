package utils

import (
	"log"
	"os"

	"github.com/bobllor/gologger"
)

// Deps is used to provide utilities to structs.
type Deps struct {
	// Log is the logging struct.
	Log *gologger.Logger
}

// NewDeps creates a new Deps.
func NewDeps(logger *gologger.Logger) *Deps {
	return &Deps{
		Log: logger,
	}
}

// NewTestDeps creates a new Deps that is preconfigured with
// fields to be used in tests.
func NewTestDeps() *Deps {
	printer := log.New(os.Stdout, "", log.Ltime|log.Ldate)
	logger := gologger.NewLogger(printer, gologger.Lsilent)

	cfg := &Deps{
		Log: logger,
	}

	return cfg
}
