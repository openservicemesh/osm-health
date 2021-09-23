package osm

import (
	"helm.sh/helm/v3/pkg/action"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm/controller"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm/pkg/k8s"
)

// ControlPlaneStatus determines the status of the OSM control plane.
func ControlPlaneStatus(osmControlPlaneNamespace common.MeshNamespace, localPort uint16, actionConfig *action.Configuration) error {
	log.Info().Msgf("Determining the status of the OSM control plane in namespace %s", osmControlPlaneNamespace)

	client, err := pod.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	controllerPods := k8s.GetOSMControllerPods(client, osmControlPlaneNamespace.String())

	outcomes := runner.Run(
		HasNoBadOsmControllerLogsCheck(client, osmControlPlaneNamespace),
		HasNoBadOsmInjectorLogsCheck(client, osmControlPlaneNamespace),
		controller.NewHTTPServerHealthEndpointsCheck(
			client,
			osmControlPlaneNamespace,
			controllerPods,
			localPort,
			actionConfig),
		controller.NewHTTPServerProxyConnectionMetricsCheck(
			client,
			osmControlPlaneNamespace,
			controllerPods,
			localPort,
			actionConfig),
	)

	printer.Print(outcomes...)

	return nil
}
