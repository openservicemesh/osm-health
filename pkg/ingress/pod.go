package ingress

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
)

// ToPod checks the Ingress to the given pod.
func ToPod(client kubernetes.Interface, toPod *v1.Pod) {
	log.Info().Msgf("Testing ingress to pod %s/%s", toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	destinationPod := k8s.Pod{
		Namespace: k8s.Namespace(toPod.Namespace),
		Name:      toPod.Name,
	}

	outcomes := common.Run(
		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, destinationPod.Namespace),
		namespace.IsMonitoredBy(client, destinationPod.Namespace, meshName),
	)

	common.Print(outcomes...)
}
