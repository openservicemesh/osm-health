package ingress

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/osm/utils"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// ToDestinationPod checks the Ingress to the given pod.
func ToDestinationPod(client kubernetes.Interface, dstPod *corev1.Pod, osmControlPlaneNamespace common.MeshNamespace) {
	log.Info().Msgf("Testing ingress to pod %s/%s", dstPod.Namespace, dstPod.Name)

	meshInfo, err := utils.GetMeshInfo(client, osmControlPlaneNamespace)
	if err != nil {
		log.Err(err).Msg("Error getting OSM info")
	}

	outcomes := runner.Run(
		// Check destination Pod's namespace
		namespace.NewSidecarInjectionCheck(client, dstPod.Namespace),
		namespace.NewMonitoredCheck(client, dstPod.Namespace, meshInfo.Name),
	)

	printer.Print(outcomes...)
}
