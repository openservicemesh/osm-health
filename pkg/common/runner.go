package common

// Run evaluates all the Runnables and returns the outcomes.
func Run(checks ...Runnable) []Outcome {
	outcomes := make([]Outcome, len(checks))
	for idx, check := range checks {
		err := check.Run()
		outcomes[idx] = Outcome{
			RunnableInfo: check.Info(),
			Error:        err,
		}
	}
	return outcomes
}
