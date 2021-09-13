package main

import (
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

func newControlPlaneCmd(actionConfig *action.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "control-plane",
		Short: "Checks the osm control plane",
		Long:  `Checks the osm control plane`,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newControlPlaneStatusCmd(actionConfig))
	return cmd
}
