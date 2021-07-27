package main

import "github.com/spf13/cobra"

func newCollectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collects various artifacts into an archive",
		Long: `
Collects various artifacts into an archive
	(add more descriptive description)
`,
		Args: cobra.NoArgs,
	}
	return cmd
}
