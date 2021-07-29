package connectivity

import (
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
	}

	outcomes := common.Run(
		// Check source Pod's namespace
		namespace.IsInjectEnabled(client, fromPod.Namespace),
		namespace.IsMonitoredBy(client, fromPod.Namespace, meshName),
		pod.HasEnvoySidecar(fromPod),

		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, toPod.Namespace),
		namespace.IsMonitoredBy(client, toPod.Namespace, meshName),
		pod.HasEnvoySidecar(toPod),
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
