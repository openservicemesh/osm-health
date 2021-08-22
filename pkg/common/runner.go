package common

import (
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Run evaluates all the Runnables and returns the outcomes.
func Run(checks ...Runnable) []Printable {
	printableOutcomes := make([]Printable, len(checks))
	for idx, check := range checks {
		outcome := check.Run()
		if outcome == nil {
			outcome = outcomes.UnknownOutcome{}
		}
		printableOutcomes[idx] = Printable{
			// TODO add check.Suggestion() and check.FixIt() in the future.
			CheckDescription: check.Description(),
			ShortStatus:      outcome.GetShortStatus(),
			LongDiagnostics:  outcome.GetLongDiagnostics(),
			Error:            outcome.GetError(),
		}
	}
	return printableOutcomes
}
