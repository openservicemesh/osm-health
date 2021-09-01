package outcomes

var _ Outcome = (*UnknownOutcome)(nil)

// UnknownOutcome is for check outcomes where the outcome isn't known.
type UnknownOutcome struct{}

// GetShortStatus implements outcomes.Outcome.
func (UnknownOutcome) GetShortStatus() string {
	return "Unknown"
}

// GetLongDiagnostics implements outcomes.Outcome.
func (o UnknownOutcome) GetLongDiagnostics() string {
	return "Unknown outcome - this check may be running into issues"
}

// GetError implements outcomes.Outcome.
func (o UnknownOutcome) GetError() error {
	return nil
}
