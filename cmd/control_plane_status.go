package main

import (
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"

	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm/pkg/constants"
)

const controlPlaneStatusExample = `$ osm-health ingress to-pod namespace-a/pod-a`

func newControlPlaneStatusCmd(actionConfig *action.Configuration) *cobra.Command {
	var localPort uint16
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Checks the status of the osm control plane",
		Example: controlPlaneStatusExample,
		Long:    `Checks the status of the osm control plane`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			osmControlPlaneNamespace := settings.Namespace()
			return osm.ControlPlaneStatus(osmControlPlaneNamespace, localPort, actionConfig)
		},
	}

	f := cmd.Flags()
	f.Uint16VarP(&localPort, "local-port", "p", constants.OSMHTTPServerPort, "Local port to use for port forwarding")

	return cmd
}
