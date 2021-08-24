package common

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Print prints the printable outcomes of the evaluation of a list of Runnables.
func Print(printables ...Printable) {
	errorsCount := 0
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 4, 0, ' ', 0)

	defer func() { _ = w.Flush() }()
	for idx, printableOutcome := range printables {
		fmt.Fprintf(w, "%d\t%s\t\t%s\n", idx+1, printableOutcome.ShortStatus, printableOutcome.CheckDescription)
		if printableOutcome.Error != nil {
			fmt.Fprintln(w, color.RedString("↳ Error: "+printableOutcome.Error.Error()))
			errorsCount = errorsCount + 1
		}
		if printableOutcome.LongDiagnostics != "" {
			fmt.Fprintln(w, "↳ Additional diagnostic info:", printableOutcome.LongDiagnostics)
		}
	}

	var statusIcon = "✅"
	if errorsCount > 0 {
		statusIcon = "❌"
	}

	fmt.Fprintf(w, "\n%s Ran %d checks. %d checks failed.\n", statusIcon, len(printables), errorsCount)
}
