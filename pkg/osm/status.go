package osm

import (
	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// ControlPlaneStatus determines the status of the OSM control plane.
func ControlPlaneStatus(osmControlPlaneNamespace common.MeshNamespace) error {
	log.Info().Msgf("Determining the status of the OSM control plane in namespace %s", osmControlPlaneNamespace)

	client, err := pod.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	outcomes := runner.Run(HasNoBadOsmControllerLogsCheck(client, osmControlPlaneNamespace))

	printer.Print(outcomes...)

	return nil
}
