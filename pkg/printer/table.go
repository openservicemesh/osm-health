package printer

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Print prints the printable outcomes of the evaluation of a list of Runnables.
func Print(printables ...common.Printable) {
	errorsCount := 0
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 4, 0, ' ', 0)

	defer func() { _ = w.Flush() }()
	for idx, printableOutcome := range printables {
		_, err := fmt.Fprintf(w, "%d\t%s\t\t%s\n", idx+1, printableOutcome.Type, printableOutcome.CheckDescription)
		if err != nil {
			return
		}
		if printableOutcome.Error != nil {
			_, err := fmt.Fprintln(w, color.RedString("---> Error: "+printableOutcome.Error.Error()))
			if err != nil {
				log.Error().Err(err)
				return
			}
			errorsCount = errorsCount + 1
		}
		if printableOutcome.Diagnostics != "" {
			_, err := fmt.Fprintln(w, "---> Diagnostic info:", printableOutcome.Diagnostics)
			if err != nil {
				log.Error().Err(err)
				return
			}
		}
	}

	var statusIcon = "✅"
	if errorsCount > 0 {
		statusIcon = "❌"
	}

	_, err := fmt.Fprintf(w, "\n%s Ran %d checks. %d checks failed.\n", statusIcon, len(printables), errorsCount)
	if err != nil {
		log.Error().Err(err)
		return
	}
}
