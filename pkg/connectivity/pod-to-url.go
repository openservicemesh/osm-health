package connectivity

import (
	"net/url"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// PodToURL tests the connectivity between a source pod and destination url.
func PodToURL(fromPod *corev1.Pod, destinationURL *url.URL) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s", fromPod.Namespace, fromPod.Name, destinationURL)

	// TODO - just like the TODO in pod-to-pod.go
	osmVersion := osm.ControllerVersion("v0.9")

	srcConfigGetter, err := envoy.GetEnvoyConfigGetterForPod(fromPod, osmVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", fromPod.Namespace, fromPod.Name)
	}

	outcomes := common.Run(
		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination URL.
		envoy.NewOutboundRouteDomainHostCheck(srcConfigGetter, destinationURL.Host),
	)

	common.Print(outcomes...)
}
