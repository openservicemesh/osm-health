package main

import (
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/osm"
)

func newControlPlaneStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Checks the status of the osm control plane",
		Example: `TODO add example`,
		Long:    `Checks the status of the osm control plane`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			osmControlPlaneNamespace := settings.Namespace()
			return osm.ControlPlaneStatus(osmControlPlaneNamespace)
		},
	}
}
