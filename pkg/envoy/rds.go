package envoy

import (
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Verify interface compliance
var _ common.Runnable = (*RouteDomainCheck)(nil)

// RouteDomainCheck implements common.Runnable
type RouteDomainCheck struct {
	*corev1.Pod
	ConfigGetter
	RouteName string
	Domain    string
}

// Run implements common.Runnable
func (l RouteDomainCheck) Run() outcomes.Outcome {
	if l.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.FailedOutcome{Error: ErrIncorrectlyInitializedConfigGetter}
	}

	envoyConfig, err := l.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.FailedOutcome{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.FailedOutcome{Error: ErrEnvoyConfigEmpty}
	}

	foundAnyRouteDomains := false
	foundSpecificRouteDomain := false

	for _, rawDynRouteCfg := range envoyConfig.Routes.GetDynamicRouteConfigs() {
		var dynRouteCfg envoy_config_route_v3.RouteConfiguration
		if err = rawDynRouteCfg.GetRouteConfig().UnmarshalTo(&dynRouteCfg); err != nil {
			return outcomes.FailedOutcome{Error: ErrUnmarshalingDynamicRouteConfig}
		}

		if dynRouteCfg.Name != l.RouteName {
			continue
		}

		for _, virtualHost := range dynRouteCfg.GetVirtualHosts() {
			for _, domain := range virtualHost.GetDomains() {
				foundAnyRouteDomains = true
				if l.Domain == "" {
					break
				}
				if domain == l.Domain {
					foundSpecificRouteDomain = true
					break
				}
			}
			if (l.Domain == "" && foundAnyRouteDomains) || foundSpecificRouteDomain {
				break
			}
		}
	}

	if !foundAnyRouteDomains {
		log.Error().Msgf("must have at least one dynamic route config domain: %+v", envoyConfig.Routes.GetDynamicRouteConfigs())
		return outcomes.FailedOutcome{Error: ErrNoDynamicRouteConfigDomains}
	}

	if l.Domain != "" && !foundSpecificRouteDomain {
		return outcomes.FailedOutcome{Error: ErrDynamicRouteConfigDomainNotFound}
	}

	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable
func (l RouteDomainCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l RouteDomainCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (l RouteDomainCheck) Description() string {
	return fmt.Sprintf("Checking whether %s is configured with correct %s Envoy route", l.ConfigGetter.GetObjectName(), l.RouteName)
}

// NewOutboundRouteDomainPodCheck creates a DestinationEndpointCheck which checks whether the Envoy config has an outbound dynamic route domain to the Pod
func NewOutboundRouteDomainPodCheck(configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		RouteName:    OutboundDynamicRouteConfigName,
		Domain:       fmt.Sprintf("%s.%s", pod.Name, pod.Namespace),
	}
}

// NewInboundRouteDomainPodCheck creates a DestinationEndpointCheck which checks whether the Envoy config has an inbound dynamic route domain to the Pod
func NewInboundRouteDomainPodCheck(configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		RouteName:    InboundDynamicRouteConfigName,
		Domain:       fmt.Sprintf("%s.%s", pod.Name, pod.Namespace),
	}
}

// NewOutboundRouteDomainHostCheck creates a DestinationEndpointCheck which checks whether the Envoy config has an outbound dynamic route domain to the URL
func NewOutboundRouteDomainHostCheck(configGetter ConfigGetter, destinationHost string) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		RouteName:    OutboundDynamicRouteConfigName,
		Domain:       destinationHost,
	}
}
