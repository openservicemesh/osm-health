package connectivity

import (
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	sourcePod := k8s.Pod{
		Namespace: k8s.Namespace(fromPod.Namespace),
		Name:      fromPod.Name,
	}

	destinationPod := k8s.Pod{
		Namespace: k8s.Namespace(toPod.Namespace),
		Name:      toPod.Name,
	}

	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
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
