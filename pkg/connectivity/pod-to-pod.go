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
func PodToPod(srcPod *corev1.Pod, dstPod *corev1.Pod, osmControlPlaneNamespace string) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", srcPod.Namespace, srcPod.Name, dstPod.Namespace, dstPod.Name)

	// TODO
	meshName := common.MeshName("osm")
	osmVersion := osm.ControllerVersion("v0.9")

	osmNamespace := common.MeshNamespace(osmControlPlaneNamespace)
	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	kubeConfig, err := kuberneteshelper.GetKubeConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error getting Kubernetes config")
	}
	splitClient, err := smiSplitClient.NewForConfig(kubeConfig)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing SMI split client")
	}

	accessClient, err := smiAccessClient.NewForConfig(kubeConfig)
	if err != nil {
		log.Err(err).Msg("Error initializing SMI access client")
	}

	var srcConfigGetter, dstConfigGetter envoy.ConfigGetter

	srcConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(srcPod, osmVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", srcPod.Namespace, srcPod.Name)
	}

	dstConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(dstPod, osmVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", dstPod.Namespace, dstPod.Name)
	}

	configurator := kuberneteshelper.GetOsmConfigurator(osmNamespace)

	outcomes := common.Run(
		// Check that pod namespaces are in the same mesh
		namespace.NewNamespacesInSameMeshCheck(client, srcPod.Namespace, dstPod.Namespace),

		// Check both pods for osm init and envoy container validity
		namespace.NewSidecarInjectionCheck(client, srcPod.Namespace),
		namespace.NewSidecarInjectionCheck(client, dstPod.Namespace),
		namespace.NewMonitoredCheck(client, srcPod.Namespace, meshName),
		namespace.NewMonitoredCheck(client, dstPod.Namespace, meshName),
		podhelper.NewMinNumContainersCheck(srcPod, 2),
		podhelper.NewMinNumContainersCheck(dstPod, 2),
		podhelper.NewOsmContainerImageCheck(configurator, srcPod),
		podhelper.NewOsmContainerImageCheck(configurator, dstPod),
		podhelper.NewEnvoySidecarImageCheck(configurator, srcPod),
		podhelper.NewEnvoySidecarImageCheck(configurator, dstPod),
		podhelper.NewProxyUUIDLabelCheck(srcPod),
		podhelper.NewProxyUUIDLabelCheck(dstPod),

		podhelper.NewEndpointsCheck(client, dstPod),

		// Check pods for bad events
		podhelper.NewPodEventsCheck(client, srcPod),
		podhelper.NewPodEventsCheck(client, dstPod),

		// Check envoy logs
		envoy.NewBadLogsCheck(client, srcPod),
		envoy.NewBadLogsCheck(client, dstPod),

		// Check osm-init logs
		osm.HasNoBadOsmInitLogsCheck(client, srcPod),
		osm.HasNoBadOsmInitLogsCheck(client, dstPod),

		// The source Envoy must have at least one endpoint for the destination Envoy.
		envoy.NewDestinationEndpointCheck(srcConfigGetter),

		// Check whether the source Pod has an endpoint that matches the destination Pod.
		envoy.NewSpecificEndpointCheck(srcConfigGetter, dstPod),

		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination Pod.
		envoy.NewOutboundRouteDomainPodCheck(srcConfigGetter, dstPod),

		// Check whether the destination Pod has an inbound dynamic route config domain that matches the source Pod.
		envoy.NewInboundRouteDomainPodCheck(dstConfigGetter, srcPod),

		// Source Envoy must have Outbound listener
		envoy.NewOutboundListenerCheck(srcConfigGetter, osmVersion),

		// Destination Envoy must have Inbound listener
		envoy.NewInboundListenerCheck(dstConfigGetter, osmVersion),

		// Source Envoy must define a cluster for the destination
		envoy.NewClusterCheck(client, srcConfigGetter, dstPod),

		// Check Envoy certificates for both pods
		envoy.HasOutboundRootCertificate(client, srcConfigGetter, dstPod),
		envoy.HasInboundRootCertificate(client, dstConfigGetter, dstPod),
		envoy.HasServiceCertificate(client, srcConfigGetter, srcPod),
		envoy.HasServiceCertificate(client, dstConfigGetter, dstPod),

		// Run SMI checks
		smi.NewTrafficSplitCheck(client, dstPod, splitClient),
		access.NewTrafficTargetCheck(osmVersion, configurator, srcPod, dstPod, accessClient),
		access.NewRoutesValidityCheck(osmVersion, configurator, srcPod, dstPod, accessClient),
	)

	common.Print(outcomes...)
}
