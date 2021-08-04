package common

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Print prints the outcomes of the evaluation of a list of Runnables.
func Print(outcomes ...Outcome) {
	issuesCount := 0
	foundIssues := false
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, ' ', 0)

	defer func() { _ = w.Flush() }()
	for idx, issue := range outcomes {
		foundIssues = foundIssues && issue.Error != nil
		errString := color.GreenString("OK")
		if issue.Error != nil {
			errString = fmt.Sprintf("%s %s", color.RedString("FAIL:"), color.RedString(issue.Error.Error()))
			issuesCount = issuesCount + 1
		}
		fmt.Fprintf(w, "%d   %s\t%s\t\n", idx+1, issue.RunnableInfo, errString)
	}

	if foundIssues {
		fmt.Printf("Found %d issues: %+v\n", issuesCount, outcomes)
	}
}
