package outcomes

// Outcome is the printable context returned from a check (common.Runnable).
type Outcome interface {
	// GetOutcomeType returns the type of the check outcome: pass/fail/info/unknown.
	GetOutcomeType() string

	// GetDiagnostics returns detailed diagnostics that were dynamically-generated during the check.
	// Diagnostics may include details about a test failure or information about the mesh/its components
	GetDiagnostics() string

	// GetError returns the error which common.Runnable{}.Run() returned.
	GetError() error
}
