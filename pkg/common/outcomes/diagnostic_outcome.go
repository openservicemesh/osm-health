package outcomes

var _ Outcome = (*DiagnosticOutcome)(nil)

// DiagnosticOutcome is the check outcome for diagnostic checks that simply give information about the status of the mesh/its components.
// These cannot be categorized as successful or failed.
// An example diagnostic check is showing whether a pod participates in a traffic split or not.
type DiagnosticOutcome struct {
	LongDiagnostics string
}

// GetShortStatus implements outcomes.Outcome.
func (DiagnosticOutcome) GetShortStatus() string {
	return "💬 Diagnostic"
}

// GetLongDiagnostics implements outcomes.Outcome.
func (o DiagnosticOutcome) GetLongDiagnostics() string {
	return o.LongDiagnostics
}

// GetError implements outcomes.Outcome.
func (o DiagnosticOutcome) GetError() error {
	return nil
}
