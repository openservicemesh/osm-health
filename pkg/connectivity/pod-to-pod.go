package connectivity

import (
	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	smiSplitClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/smi"
	"github.com/openservicemesh/osm-health/pkg/smi/access"
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

	accessClient, err := smiAccessClient.NewForConfig(kubeConfig)
	if err != nil {
		log.Err(err).Msg("Error initializing SMI access client")
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
		namespace.NewNamespacesInSameMeshCheck(client, fromPod.Namespace, toPod.Namespace),

		// Check both pods for osm init and envoy container validity
		namespace.NewSidecarInjectionCheck(client, fromPod.Namespace),
		namespace.NewSidecarInjectionCheck(client, toPod.Namespace),
		namespace.NewMonitoredCheck(client, fromPod.Namespace, meshName),
		namespace.NewMonitoredCheck(client, toPod.Namespace, meshName),
		podhelper.NewMinNumContainersCheck(fromPod, 2),
		podhelper.NewMinNumContainersCheck(toPod, 2),
		podhelper.NewOsmContainerImageCheck(configurator, fromPod),
		podhelper.NewOsmContainerImageCheck(configurator, toPod),
		podhelper.NewEnvoySidecarImageCheck(configurator, fromPod),
		podhelper.NewEnvoySidecarImageCheck(configurator, toPod),
		podhelper.NewProxyUUIDLabelCheck(fromPod),
		podhelper.NewProxyUUIDLabelCheck(toPod),

		podhelper.NewEndpointsCheck(client, toPod),

		// Check pods for bad events
		podhelper.NewPodEventsCheck(client, fromPod),
		podhelper.NewPodEventsCheck(client, toPod),

		// Check envoy logs
		envoy.NewBadLogsCheck(client, fromPod),
		envoy.NewBadLogsCheck(client, toPod),

		// The source Envoy must have at least one endpoint for the destination Envoy.
		envoy.NewDestinationEndpointCheck(srcConfigGetter),

		// Check whether the source Pod has an endpoint that matches the destination Pod.
		envoy.NewSpecificEndpointCheck(srcConfigGetter, toPod),

		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination Pod.
		envoy.NewOutboundRouteDomainPodCheck(srcConfigGetter, toPod),

		// Check whether the destination Pod has an inbound dynamic route config domain that matches the source Pod.
		envoy.NewInboundRouteDomainPodCheck(dstConfigGetter, fromPod),

		// Source Envoy must have Outbound listener
		envoy.NewOutboundListenerCheck(srcConfigGetter, osmVersion),

		// Destination Envoy must have Inbound listener
		envoy.NewInboundListenerCheck(dstConfigGetter, osmVersion),

		// Source Envoy must define a cluster for the destination
		envoy.NewClusterCheck(client, srcConfigGetter, toPod),

		// Run SMI checks
		smi.NewTrafficSplitCheck(client, toPod, splitClient),
		access.NewTrafficTargetCheck(osmVersion, configurator, fromPod, toPod, accessClient),
		access.NewRoutesValidityCheck(osmVersion, configurator, fromPod, toPod, accessClient),
	)

	common.Print(outcomes...)
}
