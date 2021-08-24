package outcomes

// Outcome is the printable context returned from a check (common.Runnable).
type Outcome interface {
	// GetShortStatus returns a short status of the check outcome, such as success, fail, or diagnostic.
	GetShortStatus() string

	// GetLongDiagnostics returns detailed diagnostics that were dynamically-generated during the check.
	// Certain checks (such as checking whether the pod participates in an SMI policy)
	// cannot be easily categorized as successful or failed but instead as diagnostic checks.
	// These diagnostic checks require dynamically-generated diagnostic info to be
	// shown to the user to give them more context.
	GetLongDiagnostics() string

	// GetError returns the error which common.Runnable{}.Run() returned.
	GetError() error
}
