package outcomes

// DependencyChecks contains checks with their outcome
var DependencyChecks = make(map[string]bool)

// GetCheckOutcome gets check outcome
func GetCheckOutcome(checkName string) bool {
	return DependencyChecks[checkName]
}

// AddCheckToMap adds the check with default outcome to false to the DependencyChecks map
func AddCheckToMap(checkName string) {
	DependencyChecks[checkName] = false
}

// SetCheckPass sets the check outcome to true if check passes
func SetCheckPass(checkName string) {
	DependencyChecks[checkName] = true
}
