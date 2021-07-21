package connectivity

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"strings"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO: actually test connectivity
	fmt.Printf("Pod %s has an envoy sidecar: %t\n", fromPod.Name, hasEnvoySideCar(fromPod))
	fmt.Printf("Pod %s has an envoy sidecar: %t\n", fromPod.Name, hasEnvoySideCar(toPod))

	return common.Result{
		SMIPolicy: common.SMIPolicy{
			HasPolicy:                  false,
			ValidPolicy:                false,
			SourceToDestinationAllowed: false,
		},
		Successful: false,
	}
}

func hasEnvoySideCar(pod *v1.Pod) bool {
	for _, container := range pod.Spec.Containers {
		if strings.Contains(container.Image, "envoy") {
			return true
		}
	}
	return false
}
