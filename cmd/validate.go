package main

import "github.com/spf13/cobra"

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates YAML files including SMI policies",
		Long: `
Validates YAML files including SMI policies
	(add more descriptive description)
`,
		Args: cobra.NoArgs,
	}
	return cmd
}
