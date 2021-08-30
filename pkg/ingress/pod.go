package ingress

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
)

// ToPod checks the Ingress to the given pod.
func ToPod(client kubernetes.Interface, toPod *corev1.Pod) {
	log.Info().Msgf("Testing ingress to pod %s/%s", toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	outcomes := common.Run(
		// Check destination Pod's namespace
		namespace.NewSidecarInjectionCheck(client, toPod.Namespace),
		namespace.NewMonitoredCheck(client, toPod.Namespace, meshName),
	)

	common.Print(outcomes...)
}
