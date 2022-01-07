package cli

import (
	"github.com/spf13/cobra"
)

// ExactArgsWithError returns the error if there are not exactly n args.
func ExactArgsWithError(n int, e error) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return e
		}
		return nil
	}
}
