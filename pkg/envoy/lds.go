package envoy

import (
	"context"
	"fmt"
	"strings"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/pkg/errors"
	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm/version"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/smi"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha2"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha3"
	"github.com/openservicemesh/osm-health/pkg/utils"
	"github.com/openservicemesh/osm/pkg/configurator"
)

// Verify interface compliance
var _ runner.Runnable = (*ListenerCheck)(nil)

// ListenerCheck implements common.Runnable
type ListenerCheck struct {
	ConfigGetter
	version.ControllerVersion

	// This is used for Info() function from Runnable. Helps the logs identify what kind of a listener we are looking for.
	listenerType string

	// This map will be used to get the expected NAME of the Envoy listener depending on the OSM version in use.
	expectedListenersPerVersion map[version.ControllerVersion]string
}

// ListenerFilterCheck implements common.Runnable
type ListenerFilterCheck struct {
	srcConfigGetter ConfigGetter
	dstConfigGetter ConfigGetter
	osmVersion      version.ControllerVersion
	cfg             configurator.Configurator
	srcPod          *corev1.Pod
	dstPod          *corev1.Pod
	accessClient    smiAccessClient.Interface
	// This is used for Info() function from Runnable. Helps the logs identify what kind of a listener we are looking for.
	listenerType string
	k8s          kubernetes.Interface
}

// FilterChainType is the prefix for the filter chain name
type FilterChainType string

const (
	// InboundMeshHTTPFilterChainPrefix is the prefix for an inbound http filter chain
	InboundMeshHTTPFilterChainPrefix FilterChainType = "inbound-mesh-http-filter-chain"

	// OutboundMeshHTTPFilterChainPrefix is the prefix for an outbound http filter chain
	OutboundMeshHTTPFilterChainPrefix FilterChainType = "outbound-mesh-http-filter-chain"

	// InboundMeshTCPFilterChainPrefix is the prefix for inbound tcp filter chain
	InboundMeshTCPFilterChainPrefix FilterChainType = "inbound-mesh-tcp-filter-chain"

	// OutboundMeshTCPFilterChainPrefix is the prefix for an outbound tcp filter chain
	OutboundMeshTCPFilterChainPrefix FilterChainType = "outbound-mesh-tcp-filter-chain"
)

// Run implements common.Runnable
func (l ListenerCheck) Run() outcomes.Outcome {
	if l.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	envoyConfig, err := l.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	expectedListenerName, exists := l.expectedListenersPerVersion[l.ControllerVersion]
	if !exists {
		return outcomes.Fail{Error: ErrOSMControllerVersionUnrecognized}
	}

	if envoyConfig == nil || envoyConfig.Listeners.DynamicListeners == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	found := false
	var actualListeners []string
	for _, actualListener := range envoyConfig.Listeners.GetDynamicListeners() {
		actualListeners = append(actualListeners, actualListener.Name)
		if expectedListenerName == actualListener.Name {
			found = true
			break
		}
	}

	if !found {
		log.Error().Msgf("must have listener with name %s but only found %s", expectedListenerName, actualListeners)
		return outcomes.Fail{Error: ErrEnvoyListenerMissing}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (l ListenerCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l ListenerCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (l ListenerCheck) Description() string {
	return fmt.Sprintf("Checking whether %s is configured with correct %s Envoy listener", l.ConfigGetter.GetObjectName(), l.listenerType)
}

// NewOutboundListenerCheck creates a ListenerCheck which checks whether the given Pod has an Envoy with properly configured listener.
func NewOutboundListenerCheck(configGetter ConfigGetter, osmVersion version.ControllerVersion) ListenerCheck {
	return ListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "outbound",

		expectedListenersPerVersion: version.OutboundListenerNames,
	}
}

// NewInboundListenerCheck creates a ListenerCheck which checks whether the given Pod has an Envoy with properly configured listener.
func NewInboundListenerCheck(configGetter ConfigGetter, osmVersion version.ControllerVersion) ListenerCheck {
	return ListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "inbound",

		expectedListenersPerVersion: version.InboundListenerNames,
	}
}

// Run implements common.Runnable
func (l ListenerFilterCheck) Run() outcomes.Outcome {
	// Check if permissive mode is enabled, in which case every meshed pod is allowed to communicate with each other
	if l.cfg.IsPermissiveTrafficPolicyMode() {
		return outcomes.Info{Diagnostics: "OSM is in permissive traffic policy modes -- all meshed pods can communicate and SMI access policies are not applicable"}
	}

	// Get rule version from TrafficTarget. The rules will be used to determine what filter chains are expected in the src and dst Envoy configs
	var ruleTypes map[string]struct{}
	var err error
	switch version.SupportedTrafficTarget[l.osmVersion] {
	case version.V1Alpha2:
		ruleTypes, err = getRuleTypesFromMatchingTrafficTargetsV1alpha2(l.srcPod, l.dstPod, l.accessClient)
	case version.V1Alpha3:
		ruleTypes, err = getRuleTypesFromMatchingTrafficTargetsV1alpha3(l.srcPod, l.dstPod, l.accessClient)
	default:
		return outcomes.Fail{Error: fmt.Errorf(
			"OSM Controller version could not be mapped to a TrafficTarget version. Supported versions are v0.5 through v0.9")}
	}
	if err != nil {
		return outcomes.Info{Diagnostics: fmt.Sprintf(
			"Pod '%s/%s' is not allowed to communicate to pod '%s/%s' via any SMI TrafficTarget policy.\n",
			l.srcPod.Namespace,
			l.srcPod.Name,
			l.dstPod.Namespace,
			l.dstPod.Name)}
	}
	if len(ruleTypes) == 0 {
		return outcomes.Info{Diagnostics: fmt.Sprintf(
			"No applicable Traffic Targets in namespace %s to check routes for, or no rules specified in Traffic Targets",
			l.dstPod.Namespace)}
	}

	// Get possible backing service(s) for dst pod
	svcs, err := pod.GetMatchingServices(l.k8s, l.dstPod.Labels, l.dstPod.Namespace)
	if err != nil {
		return outcomes.Fail{Error: errors.Wrapf(err, "failed to map Pod %s/%s to Kubernetes Services", l.dstPod.Namespace, l.dstPod.Name)}
	}

	// Get possible inbound filter chain names from the dst pod's backing service(s)
	possibleInboundFilterChainNames, err := getPossibleInboundFilterChainNames(svcs, ruleTypes)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if len(possibleInboundFilterChainNames) == 0 {
		return outcomes.Fail{Error: fmt.Errorf("failed to determine possible inbound filter chain names from dst Pod %s/%s", l.dstPod.Namespace, l.dstPod.Name)}
	}

	// Get possible outbound filter chain names from the dst pod's backing service(s)
	possibleOutboundFilterChainNames, err := getPossibleOutboundFilterChainNames(svcs, ruleTypes)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if len(possibleOutboundFilterChainNames) == 0 {
		return outcomes.Fail{Error: fmt.Errorf("failed to determine possible outbound filter chain names from dst Pod %s/%s", l.dstPod.Namespace, l.dstPod.Name)}
	}

	// Get src Envoy config
	if l.srcConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized src ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	srcEnvoyConfig, err := l.srcConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if srcEnvoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	// Get dst Envoy config
	if l.dstConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized dst ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	dstEnvoyConfig, err := l.dstConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if dstEnvoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	// Check inbound listener filter chain names in dst config
	expectedInboundListenerName, exists := version.InboundListenerNames[l.osmVersion]
	if !exists {
		return outcomes.Fail{Error: ErrOSMControllerVersionUnrecognized}
	}

	err = findMatchingFilterChainNames(dstEnvoyConfig, expectedInboundListenerName, possibleInboundFilterChainNames)
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	// Check outbound listener filter chain names in src config
	expectedOutboundListenerName, exists := version.OutboundListenerNames[l.osmVersion]
	if !exists {
		return outcomes.Fail{Error: ErrOSMControllerVersionUnrecognized}
	}

	err = findMatchingFilterChainNames(srcEnvoyConfig, expectedOutboundListenerName, possibleOutboundFilterChainNames)
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	return outcomes.Pass{}
}

func getPossibleInboundFilterChainNames(svcs []*corev1.Service, ruleTypes map[string]struct{}) (map[string]bool, error) {
	var possibleInboundFilterChainNames = map[string]bool{}
	for _, svc := range svcs {
		for ruleType := range ruleTypes {
			switch ruleType {
			case smi.HTTPRouteGroupKind:
				possibleInboundFilterChainNames[fmt.Sprintf("%s:%s", InboundMeshHTTPFilterChainPrefix, utils.K8sSvcToMeshSvc(svc).String())] = false
			case smi.TCPRouteKind:
				possibleInboundFilterChainNames[fmt.Sprintf("%s:%s", InboundMeshTCPFilterChainPrefix, utils.K8sSvcToMeshSvc(svc).String())] = false
			default:
				log.Error().Msgf("found unsupported rule type for traffic targets. Supported rule version are %s and %s. Found %s", smi.HTTPRouteGroupKind, smi.TCPRouteKind, ruleType)
				return nil, smi.ErrInvalidRuleKind
			}
		}
	}
	return possibleInboundFilterChainNames, nil
}

func getPossibleOutboundFilterChainNames(svcs []*corev1.Service, ruleTypes map[string]struct{}) (map[string]bool, error) {
	var possibleOutboundFilterChainNames = map[string]bool{}
	for _, svc := range svcs {
		for ruleType := range ruleTypes {
			switch ruleType {
			case smi.HTTPRouteGroupKind:
				possibleOutboundFilterChainNames[fmt.Sprintf("%s:%s", OutboundMeshHTTPFilterChainPrefix, utils.K8sSvcToMeshSvc(svc).String())] = false
			case smi.TCPRouteKind:
				possibleOutboundFilterChainNames[fmt.Sprintf("%s:%s", OutboundMeshTCPFilterChainPrefix, utils.K8sSvcToMeshSvc(svc).String())] = false
			default:
				log.Error().Msgf("found unsupported rule type for traffic targets. Supported rule version are %s and %s. Found %s", smi.HTTPRouteGroupKind, smi.TCPRouteKind, ruleType)
				return nil, smi.ErrInvalidRuleKind
			}
		}
	}
	return possibleOutboundFilterChainNames, nil
}

func getRuleTypesFromMatchingTrafficTargetsV1alpha2(srcPod *corev1.Pod, dstPod *corev1.Pod, accessClient smiAccessClient.Interface) (map[string]struct{}, error) {
	trafficTargets, err := accessClient.AccessV1alpha2().TrafficTargets(dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", dstPod.Namespace)
		return nil, err
	}

	ruleTypes := map[string]struct{}{}
	for _, trafficTarget := range trafficTargets.Items {
		if !v1alpha2.DoesTargetMatchPods(trafficTarget.Spec, srcPod, dstPod) {
			continue
		}
		for _, rule := range trafficTarget.Spec.Rules {
			ruleTypes[rule.Kind] = struct{}{}
		}
	}

	return ruleTypes, nil
}

func getRuleTypesFromMatchingTrafficTargetsV1alpha3(srcPod *corev1.Pod, dstPod *corev1.Pod, accessClient smiAccessClient.Interface) (map[string]struct{}, error) {
	trafficTargets, err := accessClient.AccessV1alpha3().TrafficTargets(dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", dstPod.Namespace)
		return nil, err
	}

	ruleTypes := map[string]struct{}{}
	for _, trafficTarget := range trafficTargets.Items {
		if !v1alpha3.DoesTargetMatchPods(trafficTarget.Spec, srcPod, dstPod) {
			continue
		}
		for _, rule := range trafficTarget.Spec.Rules {
			ruleTypes[rule.Kind] = struct{}{}
		}
	}

	return ruleTypes, nil
}

func findMatchingFilterChainNames(envoyConfig *Config, expectedListenerName string, possibleFilterChainNames map[string]bool) error {
	var actualFilterChainNames []string
	var actualListeners []string
	foundListener := false

	for _, actualListener := range envoyConfig.Listeners.GetDynamicListeners() {
		actualListeners = append(actualListeners, actualListener.Name)
		if expectedListenerName == actualListener.Name {
			foundListener = true
			var listener envoy_config_listener_v3.Listener
			activeStateListener := actualListener.GetActiveState().GetListener()
			if activeStateListener == nil {
				return ErrEnvoyActiveStateListenerMissing
			}
			if err := activeStateListener.UnmarshalTo(&listener); err != nil {
				return ErrUnmarshalingListener
			}
			for _, listenerFilter := range listener.FilterChains {
				log.Error().Msg(listenerFilter.Name)
				// Check filter chain name
				actualFilterChainNames = append(actualFilterChainNames, listenerFilter.Name)
				for expectedFilterChainName := range possibleFilterChainNames {
					// For osm pre v0.9, the listenerFilter.Name does not have the port number appended to it.
					// For osm v0.10 onwards, the listenerFilter.Name has the port number appended to it and may also
					// have the traffic type appended to it.
					// Examples for listenerFilter.Name v0.10 onwards:
					//				inbound-mesh-http-filter-chain:bookstore/bookstore-v1:14001
					//				outbound-mesh-http-filter-chain:bookstore/bookstore-v1_14001_http
					if strings.HasPrefix(listenerFilter.Name, expectedFilterChainName) {
						possibleFilterChainNames[expectedFilterChainName] = true
					}
				}
			}
		}
	}
	if !foundListener {
		log.Error().Msgf("must have dynamic listener with name %s but only found %v", expectedListenerName, actualListeners)
		return ErrEnvoyListenerMissing
	}
	// Iterate through map of possible filter chain names to determine if all have been found
	missingFilterChain := false
	var expectedFilterChainNames []string
	for k, v := range possibleFilterChainNames {
		expectedFilterChainNames = append(expectedFilterChainNames, k)
		if !v {
			missingFilterChain = true
		}
	}
	if missingFilterChain {
		log.Error().Msgf("expected filter chain names %v but only found %v", expectedFilterChainNames, actualFilterChainNames)
		return ErrEnvoyFilterChainMissing
	}
	return nil
}

// Suggestion implements common.Runnable
func (l ListenerFilterCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l ListenerFilterCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (l ListenerFilterCheck) Description() string {
	return fmt.Sprintf("Checking whether %s and %s are configured with the correct %s Envoy filter chains",
		l.srcConfigGetter.GetObjectName(), l.dstConfigGetter.GetObjectName(), l.listenerType)
}

// NewListenerFilterCheck creates a ListenerFilterCheck which checks whether the given Pods have Envoys with properly configured listener filter chains.
func NewListenerFilterCheck(
	srcConfigGetter ConfigGetter,
	dstConfigGetter ConfigGetter,
	osmVersion version.ControllerVersion,
	cfg configurator.Configurator,
	srcPod *corev1.Pod,
	dstPod *corev1.Pod,
	accessClient smiAccessClient.Interface,
	k8s kubernetes.Interface) ListenerFilterCheck {
	return ListenerFilterCheck{
		srcConfigGetter: srcConfigGetter,
		dstConfigGetter: dstConfigGetter,
		osmVersion:      osmVersion,
		cfg:             cfg,
		srcPod:          srcPod,
		dstPod:          dstPod,
		accessClient:    accessClient,
		k8s:             k8s,
	}
}
