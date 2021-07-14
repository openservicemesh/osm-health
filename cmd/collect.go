package main

import (
	"io"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

const collectDesc = `
Collects various artifacts into an archive
	(add more descriptive description)
`

func newCollectCmd(config *action.Configuration, in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collects various artifacts into an archive",
		Long:  collectDesc,
		Args:  cobra.NoArgs,
	}
	return cmd
}
