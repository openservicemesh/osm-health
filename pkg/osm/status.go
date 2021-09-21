package osm

import (
	"helm.sh/helm/v3/pkg/action"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm/controller"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// ControlPlaneStatus determines the status of the OSM control plane.
func ControlPlaneStatus(osmControlPlaneNamespace common.MeshNamespace, localPort uint16, actionConfig *action.Configuration) error {
	log.Info().Msgf("Determining the status of the OSM control plane in namespace %s", osmControlPlaneNamespace)

	client, err := pod.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	outcomes := runner.Run(
		HasNoBadOsmControllerLogsCheck(client, osmControlPlaneNamespace),
		controller.HasValidInfoFromControllerHTTPServerEndpointsCheck(client, osmControlPlaneNamespace, localPort, actionConfig),
	)

	printer.Print(outcomes...)

	return nil
}
