package connectivity

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(client kubernetes.Interface, fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	sourcePod := k8s.Pod{
		Namespace: k8s.Namespace(fromPod.Namespace),
		Name:      fromPod.Name,
	}

	destinationPod := k8s.Pod{
		Namespace: k8s.Namespace(toPod.Namespace),
		Name:      fromPod.Name,
	}

	outcomes := common.Run(
		// Check source Pod's namespace
		namespace.IsInjectEnabled(client, sourcePod.Namespace),
		namespace.IsMonitoredBy(client, sourcePod.Namespace, meshName),

		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, destinationPod.Namespace),
		namespace.IsMonitoredBy(client, destinationPod.Namespace, meshName),
	)

	common.Print(outcomes...)

	return common.Result{
		SMIPolicy: common.SMIPolicy{
			HasPolicy:                  false,
			ValidPolicy:                false,
			SourceToDestinationAllowed: false,
		},
		Successful: false,
	}
}
