package outcomes

var _ Outcome = (*Unknown)(nil)

// Unknown is the outcome type that occurs when the check did not have a conclusive result
type Unknown struct{}

// GetOutcomeType implements outcomes.Outcome.
func (Unknown) GetOutcomeType() string {
	return "Unknown"
}

// GetDiagnostics implements outcomes.Outcome.
func (o Unknown) GetDiagnostics() string {
	return "Unknown outcome - this check may be running into issues"
}

// GetError implements outcomes.Outcome.
func (o Unknown) GetError() error {
	return nil
}
