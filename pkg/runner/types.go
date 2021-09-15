package runner

import (
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Runnable is a type of generic function that can be executed; it returns pass/fail on a given check.
type Runnable interface {
	// Run executes a check and returns an outcome (outcomes.Outcome).
	Run() outcomes.Outcome

	// Description returns human-readable information on what check is being executed.
	Description() string

	// Suggestion returns a human-readable suggestion on how to fix the issue.
	Suggestion() string

	// FixIt attempts to fix the issue at hand.
	FixIt() error
}
