package connectivity

import (
	"net/url"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// PodToURL tests the connectivity between a source pod and destination url.
func PodToURL(srcPod *corev1.Pod, destinationURL *url.URL, osmControlPlaneNamespace common.MeshNamespace) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s", srcPod.Namespace, srcPod.Name, destinationURL)

	client, err := pod.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	meshInfo, err := osm.GetMeshInfo(client, osmControlPlaneNamespace)
	if err != nil {
		log.Error().Err(err).Msg("Error getting OSM info")
	}

	srcConfigGetter, err := envoy.GetEnvoyConfigGetterForPod(srcPod, meshInfo.OSMVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", srcPod.Namespace, srcPod.Name)
	}

	outcomes := runner.Run(
		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination URL.
		envoy.NewOutboundRouteDomainHostCheck(srcConfigGetter, destinationURL.Host),
	)

	printer.Print(outcomes...)
}
