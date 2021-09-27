package connectivity

import (
	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	smiSpecClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/specs/clientset/versioned"
	smiSplitClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/envoy"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/namespace"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/osm/utils"
	"github.com/openservicemesh/osm-health/pkg/printer"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/smi/access"
	"github.com/openservicemesh/osm-health/pkg/smi/split"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(srcPod *corev1.Pod, dstPod *corev1.Pod, osmControlPlaneNamespace common.MeshNamespace) {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", srcPod.Namespace, srcPod.Name, dstPod.Namespace, dstPod.Name)

	client, err := pod.GetKubeClient()
	if err != nil {
		log.Error().Err(err).Msg("Error creating Kubernetes client")
	}

	meshInfo, err := utils.GetMeshInfo(client, osmControlPlaneNamespace)
	if err != nil {
		log.Err(err).Msg("Error getting OSM info")
	}

	kubeConfig, err := pod.GetKubeConfig()
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

	specClient, err := smiSpecClient.NewForConfig(kubeConfig)
	if err != nil {
		log.Err(err).Msg("Error initializing SMI spec client")
	}

	var srcConfigGetter, dstConfigGetter envoy.ConfigGetter

	srcConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(srcPod, meshInfo.OSMVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", srcPod.Namespace, srcPod.Name)
	}

	dstConfigGetter, err = envoy.GetEnvoyConfigGetterForPod(dstPod, meshInfo.OSMVersion)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating ConfigGetter for pod %s/%s", dstPod.Namespace, dstPod.Name)
	}

	configurator := pod.GetOsmConfigurator(meshInfo.Namespace)

	checks := []runner.Runnable{
		// Check that pod namespaces are in the same mesh
		namespace.NewNamespacesInSameMeshCheck(client, srcPod.Namespace, dstPod.Namespace),

		// Check both pods for osm init and envoy container validity
		namespace.NewSidecarInjectionCheck(client, srcPod.Namespace),
		namespace.NewSidecarInjectionCheck(client, dstPod.Namespace),
		namespace.NewMonitoredCheck(client, srcPod.Namespace, meshInfo.Name),
		namespace.NewMonitoredCheck(client, dstPod.Namespace, meshInfo.Name),
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
		podhelper.HasNoBadOsmInitLogsCheck(client, srcPod),
		podhelper.HasNoBadOsmInitLogsCheck(client, dstPod),

		// The destination pod must have at least one service.
		podhelper.NewServiceCheck(client, dstPod),

		// The source Envoy must have at least one endpoint for the destination Envoy.
		envoy.NewDestinationEndpointCheck(srcConfigGetter),

		// Check whether the source Pod has an endpoint that matches the destination Pod.
		envoy.NewSpecificEndpointCheck(srcConfigGetter, dstPod),

		// Check whether the source Pod has an outbound dynamic route config domain that matches the destination Pod.
		envoy.NewOutboundRouteDomainPodCheck(client, srcConfigGetter, dstPod),

		// Check whether the destination Pod has an inbound dynamic route config domain that matches the source Pod.
		envoy.NewInboundRouteDomainPodCheck(client, dstConfigGetter, srcPod),

		// Source Envoy must have Outbound listener
		envoy.NewOutboundListenerCheck(srcConfigGetter, meshInfo.OSMVersion),

		// Destination Envoy must have Inbound listener
		envoy.NewInboundListenerCheck(dstConfigGetter, meshInfo.OSMVersion),

		// Source Envoy must define a cluster for the destination
		envoy.NewClusterCheck(client, srcConfigGetter, dstPod),

		// Check Envoy certificates for both pods
		envoy.HasOutboundRootCertificate(client, srcConfigGetter, dstPod),
		envoy.HasInboundRootCertificate(client, dstConfigGetter, dstPod),
		envoy.HasServiceCertificate(client, srcConfigGetter, srcPod),
		envoy.HasServiceCertificate(client, dstConfigGetter, dstPod),

		// Check Envoy for dynamic warming issues
		envoy.NewDynamicWarmingCheck(srcConfigGetter),
		envoy.NewDynamicWarmingCheck(dstConfigGetter),

		// Run SMI checks
		split.NewTrafficSplitCheck(meshInfo.OSMVersion, client, dstPod, splitClient),
		access.NewTrafficTargetCheck(meshInfo.OSMVersion, configurator, srcPod, dstPod, accessClient),
		access.NewRoutesValidityCheck(meshInfo.OSMVersion, configurator, srcPod, dstPod, accessClient),
		access.NewRoutesExistenceCheck(meshInfo.OSMVersion, configurator, srcPod, dstPod, accessClient, specClient),

		// Check whether the source and destination envoys have filter chains that match the destination service.
		envoy.NewListenerFilterCheck(srcConfigGetter, dstConfigGetter, meshInfo.OSMVersion, configurator, srcPod, dstPod, accessClient, client),
	}

	outcomes := runner.Run(checks...)
	printer.Print(outcomes...)
}
