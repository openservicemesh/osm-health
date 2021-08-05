package osm

import (
	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// ControlPlaneStatus determines the status of the OSM control plane.
func ControlPlaneStatus(osmControlPlaneNamespace string) error {
	log.Info().Msgf("Determining the status of the OSM control plane in namespace %s", osmControlPlaneNamespace)

	_, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
	}

	// TODO add checks like osm controller log checks
	outcomes := common.Run()

	common.Print(outcomes...)

	return nil
}
