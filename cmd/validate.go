package main

import (
	"io"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

const validateDesc = `
Validates YAML files including SMI policies
	(add more descriptive description)
`

func newValidateCmd(config *action.Configuration, in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates YAML files including SMI policies",
		Long:  validateDesc,
		Args:  cobra.NoArgs,
	}
	return cmd
}
