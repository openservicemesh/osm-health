package main

import (
	"io"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

const connectivityDesc = `
Checks connectivity between Kubernetes resources
	(add more descriptive description)
`

func newConnectivityCmd(config *action.Configuration, in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectivity",
		Short: "Checks connectivity between Kubernetes resources",
		Long:  connectivityDesc,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newConnectivityPodToPodCmd(config, in, out))
	return cmd
}
