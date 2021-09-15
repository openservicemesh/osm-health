package ingress

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// ToPod checks the Ingress to the given pod.
func ToPod(client kubernetes.Interface, toPod *corev1.Pod) {
	log.Info().Msgf("Testing ingress to pod %s/%s", toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")

	outcomes := runner.Run(
		// Check destination Pod's namespace
		namespace.NewSidecarInjectionCheck(client, toPod.Namespace),
		namespace.NewMonitoredCheck(client, toPod.Namespace, meshName),
	)

	printer.Print(outcomes...)
}
