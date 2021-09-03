package outcomes

import (
	"github.com/fatih/color"
)

var _ Outcome = (*Fail)(nil)

// Fail is the check outcome for checks that fail or encounter errors.
type Fail struct {
	Error error
}

// GetOutcomeType implements outcomes.Outcome.
func (Fail) GetOutcomeType() string {
	return color.RedString("Fail")
}

// GetDiagnostics implements outcomes.Outcome.
func (o Fail) GetDiagnostics() string {
	// No diagnostics information is returned as the Error field
	// should contain all the necessary information about the error.
	return NoDiagnosticInfo
}

// GetError implements outcomes.Outcome.
func (o Fail) GetError() error {
	return o.Error
}
