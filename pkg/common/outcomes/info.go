package outcomes

import (
	"github.com/fatih/color"
)

var _ Outcome = (*Info)(nil)

// Info is the check outcome for checks that cannot be categorized as pass or fail, but instead simply give information
// about the status of the mesh/its components. These checks show dynamically-generated diagnostic information to the user to provide more context.
// Ex: check whether a pod participates in an SMI TrafficSplit or not, if yes - output the name of the TrafficSplit
type Info struct {
	Diagnostics string
}

// GetOutcomeType implements outcomes.Outcome.
func (Info) GetOutcomeType() string {
	return color.BlueString("Info")
}

// GetDiagnostics implements outcomes.Outcome.
func (o Info) GetDiagnostics() string {
	return o.Diagnostics
}

// GetError implements outcomes.Outcome.
func (o Info) GetError() error {
	return nil
}
