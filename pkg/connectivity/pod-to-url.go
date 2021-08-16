package connectivity

import (
	"net/url"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodToURL tests the connectivity between a source pod and destination url.
func PodToURL(fromPod *corev1.Pod, destinationURL *url.URL) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s", fromPod.Namespace, fromPod.Name, destinationURL)
	outcomes := common.Run()
	common.Print(outcomes...)
}
