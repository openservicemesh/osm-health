package outcomes

var _ Outcome = (*FailedOutcome)(nil)

// FailedOutcome is the check outcome for checks that fail or encounter errors.
type FailedOutcome struct {
	Error error
}

// GetShortStatus implements outcomes.Outcome.
func (FailedOutcome) GetShortStatus() string {
	return "❌ Fail"
}

// GetLongDiagnostics implements outcomes.Outcome.
func (o FailedOutcome) GetLongDiagnostics() string {
	// No diagnostics information is returned as the Error field
	// should contain all the necessary information about the error.
	return NoDiagnosticInfo
}

// GetError imeplements outcomes.Outcome.
func (o FailedOutcome) GetError() error {
	return o.Error
}
