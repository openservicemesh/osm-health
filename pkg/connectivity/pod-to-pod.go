package connectivity

import (
	smiSplitClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/smi"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(fromPod *corev1.Pod, toPod *corev1.Pod, osmControlPlaneNamespace string) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", fromPod.Namespace, fromPod.Name, toPod.Namespace, toPod.Name)

	// TODO
	meshName := common.MeshName("osm")
	osmVersion := osm.ControllerVersion("v0.9")

	osmNamespace := common.MeshNamespace(osmControlPlaneNamespace)
	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Err(err).Msg("Error creating Kubernetes client")
	}

	kubeConfig, err := kuberneteshelper.GetKubeConfig()
	if err != nil {
		log.Err(err).Msg("Error getting Kubernetes config")
	}
	splitClient, err := smiSplitClient.NewForConfig(kubeConfig)
	if err != nil {
		log.Err(err).Msg("Error initializing SMI split client")
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

	configurator := kuberneteshelper.GetOsmConfigurator(osmNamespace)

	outcomes := common.Run(
		// Check that pod namespaces are in the same mesh
		namespace.AreNamespacesInSameMesh(client, fromPod.Namespace, toPod.Namespace),

		// Check source Pod's namespace
		namespace.IsInjectEnabled(client, fromPod.Namespace),
		namespace.IsMonitoredBy(client, fromPod.Namespace, meshName),
		podhelper.HasMinExpectedContainers(fromPod, 2),
		podhelper.HasExpectedOsmInitImage(configurator, fromPod),
		podhelper.HasExpectedEnvoyImage(configurator, fromPod),
		podhelper.HasProxyUUIDLabel(fromPod),
		podhelper.DoesNotHaveBadEvents(client, fromPod),

		// Check destination Pod's namespace
		namespace.IsInjectEnabled(client, toPod.Namespace),
		namespace.IsMonitoredBy(client, toPod.Namespace, meshName),
		podhelper.HasMinExpectedContainers(toPod, 2),
		podhelper.HasExpectedOsmInitImage(configurator, toPod),
		podhelper.HasExpectedEnvoyImage(configurator, toPod),
		podhelper.HasProxyUUIDLabel(toPod),
		podhelper.DoesNotHaveBadEvents(client, toPod),

		// The source Envoy must have at least one endpoint for the destination Envoy.
		envoy.HasDestinationEndpoints(srcConfigGetter),

		// Check whether the source Pod has an endpoint that matches the destination Pod.
		envoy.HasSpecificEndpoint(srcConfigGetter, toPod),

		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination Pod.
		envoy.HasOutboundDynamicRouteConfigDomainCheck(srcConfigGetter, toPod),

		// Check whether the destination Pod has an inbound dynamic route config domain that matches the source Pod.
		envoy.HasInboundDynamicRouteConfigDomainCheck(dstConfigGetter, fromPod),

		// Check envoy logs
		envoy.HasNoBadEnvoyLogsCheck(client, fromPod),
		envoy.HasNoBadEnvoyLogsCheck(client, toPod),

		// Source Envoy must have Outbound listener
		envoy.HasOutboundListener(srcConfigGetter, osmVersion),

		// Destination Envoy must have Inbound listener
		envoy.HasInboundListener(dstConfigGetter, osmVersion),

		// Source Envoy must define a cluster for the destination
		envoy.HasCluster(client, srcConfigGetter, toPod),

		// Run SMI checks
		smi.IsInTrafficSplit(client, toPod, splitClient),
	)

	common.Print(outcomes...)
}
