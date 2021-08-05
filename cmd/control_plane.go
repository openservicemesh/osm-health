package main

import "github.com/spf13/cobra"

func newControlPlaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "control-plane",
		Short: "Checks the osm control plane",
		Long:  `Checks the osm control plane`,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newControlPlaneStatusCmd())
	return cmd
}
