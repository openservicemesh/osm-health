package osm

import (
	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// ControlPlaneStatus determines the status of the OSM control plane.
func ControlPlaneStatus(osmControlPlaneNamespace string) error {
	log.Info().Msgf("Determining the status of the OSM control plane in namespace %s", osmControlPlaneNamespace)

	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	outcomes := common.Run(HasNoBadOsmControllerLogsCheck(client, osmControlPlaneNamespace))

	common.Print(outcomes...)

	return nil
}
