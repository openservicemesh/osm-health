package pod

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned"
	"github.com/openservicemesh/osm/pkg/signals"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/logger"
)

var (
	log = logger.New("osm-health/kubernetes/pod")
)

// EnvoySidecarCheck implements common.Runnable
type EnvoySidecarCheck struct {
	pod *v1.Pod
}

// HasEnvoySidecar checks whether a pod has a sidecar with the envoy image specified in the meshconfig
func HasEnvoySidecar(pod *v1.Pod) common.Runnable {
	return hasEnvoySidecar(pod)
}

func hasEnvoySidecar(pod *v1.Pod) EnvoySidecarCheck {
	return EnvoySidecarCheck{
		pod: pod,
	}
}

// Info implements common.Runnable
func (check EnvoySidecarCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has a container with envoy image matching meshconfig envoy image", check.pod.Name)
}

// Run implements common.Runnable
func (check EnvoySidecarCheck) Run() error {
	stop := signals.RegisterExitHandlers()
	kubeConfig, err := kuberneteshelper.GetKubeConfig()
	if err != nil {
		log.Err(err).Msg("Error getting kubeconfig")
	}
	cfg := configurator.NewConfigurator(versioned.NewForConfigOrDie(kubeConfig), stop, check.pod.Namespace, check.pod.Name)

	for _, container := range check.pod.Spec.Containers {
		if container.Image == cfg.GetEnvoyImage() {
			return nil
		}
	}
	return ErrExpectedEnvoySidcarMissing
}
