package connectivity

import (
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *v1.Pod, toPod *v1.Pod) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")
	osmVersion := osm.ControllerVersion("v0.9")

	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
	}

	var srcConfigGetter, dstConfigGetter envoy.ConfigGetter

	srcConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(fromPod, osmVersion)
	if err != nil {
		log.Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", fromPod.Namespace, fromPod.Name)
	}

	dstConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(toPod, osmVersion)
	if err != nil {
		log.Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", toPod.Namespace, toPod.Name)
	}

	outcomes := common.Run(
		// Check source Pod's namespace
		namespace.IsInjectEnabled(client, fromPod.Namespace),
		namespace.IsMonitoredBy(client, fromPod.Namespace, meshName),
		pod.HasMinExpectedContainers(fromPod, 3),
		pod.HasExpectedEnvoyImage(fromPod),
		pod.HasProxyUUIDLabel(fromPod),
		pod.DoesNotHaveBadEvents(client, fromPod),

		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, toPod.Namespace),
		namespace.IsMonitoredBy(client, toPod.Namespace, meshName),
		pod.HasMinExpectedContainers(toPod, 3),
		pod.HasExpectedEnvoyImage(toPod),
		pod.HasProxyUUIDLabel(toPod),
		pod.DoesNotHaveBadEvents(client, toPod),

		// The source Envoy must have at least one endpoint for the destination Envoy.
		envoy.HasDestinationEndpoints(srcConfigGetter),

		// Check whether the source Pod has an endpoint that matches the destination Pod.
		envoy.HasSpecificEndpoint(srcConfigGetter, toPod),

		// Check envoy logs
		envoy.HasNoBadEnvoyLogsCheck(client, fromPod),
		envoy.HasNoBadEnvoyLogsCheck(client, toPod),

		// Source Envoy must have Outbound listener
		envoy.HasOutboundListener(srcConfigGetter, osmVersion),

		// Destination Envoy must have Inbound listener
		envoy.HasInboundListener(dstConfigGetter, osmVersion),

		// Source Envoy must define a cluster for the destination
		envoy.HasCluster(client, srcConfigGetter, toPod),
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
