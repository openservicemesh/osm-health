package connectivity

import (
	"fmt"
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO: actually test connectivity
	fmt.Printf("Pod %s has an envoy sidecar: %t\n", fromPod.Name, hasEnvoySideCar(fromPod))
	fmt.Printf("Pod %s has an envoy sidecar: %t\n", toPod.Name, hasEnvoySideCar(toPod))

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
	foundEnvoyContainer := checkForContainer(pod.Spec.Containers, "envoy")
	foundOsmInitContainer := checkForContainer(pod.Spec.InitContainers, "osm-init")
	return foundEnvoyContainer && foundOsmInitContainer && isPodMeshed(pod) && (countNumContainers(pod) >= 3)
}

func isPodMeshed(pod *v1.Pod) bool {
	_, labels := pod.Labels["osm-proxy-uuid"] // TODO: change this to constants.EnvoyUniqueIDLabelName?
	return labels
}

func countNumContainers(pod *v1.Pod) int {
	return len(pod.Spec.Containers) + len(pod.Spec.InitContainers)
}

// TODO: more specific check?
func checkForContainer(containerList []v1.Container, name string) bool {
	for _, container := range containerList {
		if container.Name == name {
			return true
		}
	}
	return false
}