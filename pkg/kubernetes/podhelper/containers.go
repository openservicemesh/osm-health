package podhelper

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm/pkg/configurator"
)

// Verify interface compliance
var _ common.Runnable = (*EnvoySidecarImageCheck)(nil)

// EnvoySidecarImageCheck implements common.Runnable
type EnvoySidecarImageCheck struct {
	cfg configurator.Configurator
	pod *corev1.Pod
}

// HasExpectedEnvoyImage checks whether a pod has a sidecar with the envoy image specified in the meshconfig
func HasExpectedEnvoyImage(osmConfigurator configurator.Configurator, pod *corev1.Pod) EnvoySidecarImageCheck {
	return EnvoySidecarImageCheck{
		cfg: osmConfigurator,
		pod: pod,
	}
}

// Description implements common.Runnable
func (check EnvoySidecarImageCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has a container with envoy image matching meshconfig envoy image", check.pod.Name)
}

// Run implements common.Runnable
func (check EnvoySidecarImageCheck) Run() outcomes.Outcome {
	for _, container := range check.pod.Spec.Containers {
		if container.Image == check.cfg.GetEnvoyImage() {
			return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
		}
	}
	return outcomes.FailedOutcome{Error: ErrExpectedEnvoyImageMissing}
}

// Suggestion implements common.Runnable
func (check EnvoySidecarImageCheck) Suggestion() string {
	return fmt.Sprintf("Envoy image may exist but differs from image in mesh config. To verify that the pod has an envoy container, try: \"kubectl describe pod %s -n %s\"", check.pod.Name, check.pod.Namespace)
}

// FixIt implements common.Runnable
func (check EnvoySidecarImageCheck) FixIt() error {
	panic("implement me")
}

// Verify interface compliance
var _ common.Runnable = (*OsmInitContainerImageCheck)(nil)

// OsmInitContainerImageCheck implements common.Runnable
type OsmInitContainerImageCheck struct {
	cfg configurator.Configurator
	pod *corev1.Pod
}

// HasExpectedOsmInitImage checks whether a pod has a sidecar with the osm init container image specified in the meshconfig
func HasExpectedOsmInitImage(osmConfigurator configurator.Configurator, pod *corev1.Pod) OsmInitContainerImageCheck {
	return OsmInitContainerImageCheck{
		cfg: osmConfigurator,
		pod: pod,
	}
}

// Description implements common.Runnable
func (check OsmInitContainerImageCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has a container with osm init image matching meshconfig init container image", check.pod.Name)
}

// Run implements common.Runnable
func (check OsmInitContainerImageCheck) Run() outcomes.Outcome {
	for _, container := range check.pod.Spec.InitContainers {
		if container.Image == check.cfg.GetInitContainerImage() {
			return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
		}
	}
	return outcomes.FailedOutcome{Error: ErrExpectedOsmInitImageMissing}
}

// Suggestion implements common.Runnable
func (check OsmInitContainerImageCheck) Suggestion() string {
	return fmt.Sprintf("OSM init container image may exist but differs from image in mesh config. To inspect the pod's init container, try: \"kubectl describe pod %s -n %s\"", check.pod.Name, check.pod.Namespace)
}

// FixIt implements common.Runnable
func (check OsmInitContainerImageCheck) FixIt() error {
	panic("implement me")
}

// Verify interface compliance
var _ common.Runnable = (*MinNumContainersCheck)(nil)

// MinNumContainersCheck implements common.Runnable
type MinNumContainersCheck struct {
	pod    *corev1.Pod
	minNum int
}

// HasMinExpectedContainers checks whether a pod has at least the min number of containers expected
// This currently corresponds to an app container, osm init container and envoy proxy sidecar
func HasMinExpectedContainers(pod *corev1.Pod, num int) MinNumContainersCheck {
	return MinNumContainersCheck{
		pod:    pod,
		minNum: num,
	}
}

// Description implements common.Runnable
func (check MinNumContainersCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has at least %d containers", check.pod.Name, check.minNum)
}

// Run implements common.Runnable
func (check MinNumContainersCheck) Run() outcomes.Outcome {
	if len(check.pod.Spec.Containers) < check.minNum {
		return outcomes.FailedOutcome{Error: ErrExpectedMinNumContainers}
	}
	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable
func (check MinNumContainersCheck) Suggestion() string {
	return fmt.Sprintf("Verify that the pod has at least %d containers (app container and envoy sidecar container). Try: \"kubectl describe pod %s -n %s\"", check.minNum, check.pod.Name, check.pod.Namespace)
}

// FixIt implements common.Runnable
func (check MinNumContainersCheck) FixIt() error {
	panic("implement me")
}
