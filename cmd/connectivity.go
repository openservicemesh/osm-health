package main

import (
	"github.com/spf13/cobra"
)

const connectivityDesc = `
Checks connectivity between Kubernetes resources
	(add more descriptive description)
`

func newConnectivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectivity",
		Short: "Checks connectivity between Kubernetes resources",
		Long:  connectivityDesc,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newConnectivityPodToPodCmd())
	cmd.AddCommand(newConnectivityPodToURLCmd())
	return cmd
}
