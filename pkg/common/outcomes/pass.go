package outcomes

import (
	"github.com/fatih/color"
)

var _ Outcome = (*Pass)(nil)

// Pass is for check outcomes that are successful and do not have diagnostic information to show.
type Pass struct {
	Msg string
}

// GetOutcomeType implements outcomes.Outcome.
func (Pass) GetOutcomeType() string {
	return color.GreenString("Pass")
}

// GetDiagnostics implements outcomes.Outcome.
func (o Pass) GetDiagnostics() string {
	return o.Msg
}

// GetError implements outcomes.Outcome.
func (o Pass) GetError() error {
	return nil
}
