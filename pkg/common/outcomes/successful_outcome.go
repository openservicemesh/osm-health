package outcomes

var _ Outcome = (*SuccessfulOutcomeWithoutDiagnostics)(nil)

// SuccessfulOutcomeWithoutDiagnostics is for check outcomes that are successful and do not have diagnostic information to show.
type SuccessfulOutcomeWithoutDiagnostics struct{}

// GetShortStatus implements outcomes.Outcome.
func (SuccessfulOutcomeWithoutDiagnostics) GetShortStatus() string {
	return "Pass"
}

// GetLongDiagnostics implements outcomes.Outcome.
func (o SuccessfulOutcomeWithoutDiagnostics) GetLongDiagnostics() string {
	return NoDiagnosticInfo
}

// GetError implements outcomes.Outcome.
func (o SuccessfulOutcomeWithoutDiagnostics) GetError() error {
	return nil
}
