// Package logger implements utility routines to initialize the logging facility used by OSM components.
package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CallerHook implements zerolog.Hook interface.
type CallerHook struct{}

// Run adds additional context
func (h CallerHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	if _, file, line, ok := runtime.Caller(3); ok {
		e.Str("file", fmt.Sprintf("%s:%d", path.Base(file), line))
	}
}

// New creates a new zerolog.Logger
func New(component string) zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	return newLogger(component).Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func newLogger(module string) zerolog.Logger {
	return log.With().Str("module", module).Logger().Hook(CallerHook{})
}
