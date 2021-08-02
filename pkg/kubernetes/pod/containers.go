package pod

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// EnvoySidecarImageCheck implements common.Runnable
type EnvoySidecarImageCheck struct {
	pod *v1.Pod
}

// HasExpectedEnvoyImage checks whether a pod has a sidecar with the envoy image specified in the meshconfig
func HasExpectedEnvoyImage(pod *v1.Pod) common.Runnable {
	return EnvoySidecarImageCheck{
		pod: pod,
	}
}

// Info implements common.Runnable
func (check EnvoySidecarImageCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has a container with envoy image matching meshconfig envoy image", check.pod.Name)
}

// Run implements common.Runnable
func (check EnvoySidecarImageCheck) Run() error {
	cfg := kuberneteshelper.GetOsmConfigurator(check.pod)

	for _, container := range check.pod.Spec.Containers {
		if container.Image == cfg.GetEnvoyImage() {
			return nil
		}
	}
	return ErrExpectedEnvoyImageMissing
}

// MinNumContainersCheck implements common.Runnable
type MinNumContainersCheck struct {
	pod    *v1.Pod
	minNum int
}

// HasMinExpectedContainers checks whether a pod has at least the min number of containers expected
// This currently corresponds to an app container, osm init container and envoy proxy sidecar
func HasMinExpectedContainers(pod *v1.Pod, num int) common.Runnable {
	return MinNumContainersCheck{
		pod:    pod,
		minNum: num,
	}
}

// Info implements common.Runnable
func (check MinNumContainersCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has at least %d containers", check.pod.Name, check.minNum)
}

// Run implements common.Runnable
func (check MinNumContainersCheck) Run() error {
	if len(check.pod.Spec.Containers)+len(check.pod.Spec.InitContainers) < check.minNum {
		return ErrExpectedMinNumContainers
	}
	return nil
}
