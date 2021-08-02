package common

import "fmt"

// Print prints the outcomes of the evaluation of a list of Runnables.
func Print(outcomes ...Outcome) {
	issuesCount := 0
	foundIssues := false
	for idx, issue := range outcomes {
		foundIssues = foundIssues && issue.Error != nil
		errString := "OK"
		if issue.Error != nil {
			errString = fmt.Sprintf("FAIL: %s", issue.Error.Error())
			issuesCount = issuesCount + 1
		}
		fmt.Printf("%d  %s  -- %s\n", idx+1, issue.RunnableInfo, errString)
	}

	if foundIssues {
		fmt.Printf("Found %d issues: %+v\n", issuesCount, outcomes)
	}
}
