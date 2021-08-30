package outcomes

var _ Outcome = (*SkipOutcome)(nil)

// SkipOutcome is the check outcome for checks that fail or encounter errors.
type SkipOutcome struct {
	LongDiagnostics string
}

// GetShortStatus implements outcomes.Outcome.
func (SkipOutcome) GetShortStatus() string {
	return "➡️ Skipped"
}

// GetLongDiagnostics implements outcomes.Outcome.
func (o SkipOutcome) GetLongDiagnostics() string {
	return o.LongDiagnostics
}

// GetError imeplements outcomes.Outcome.
func (o SkipOutcome) GetError() error {
	return nil
}
