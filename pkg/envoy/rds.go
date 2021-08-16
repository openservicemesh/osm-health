package envoy

import (
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*RouteDomainCheck)(nil)

// RouteDomainCheck implements common.Runnable
type RouteDomainCheck struct {
	*corev1.Pod
	ConfigGetter
	RouteName string
}

// Run implements common.Runnable
func (l RouteDomainCheck) Run() error {
	if l.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return ErrIncorrectlyInitializedConfigGetter
	}

	envoyConfig, err := l.ConfigGetter.GetConfig()
	if err != nil {
		return err
	}

	if envoyConfig == nil {
		return ErrEnvoyConfigEmpty
	}

	foundAnyRouteDomains := false
	foundSpecificRouteDomain := false

	for _, rawDynRouteCfg := range envoyConfig.Routes.GetDynamicRouteConfigs() {
		var dynRouteCfg envoy_config_route_v3.RouteConfiguration
		if err = rawDynRouteCfg.GetRouteConfig().UnmarshalTo(&dynRouteCfg); err != nil {
			return ErrUnmarshalingDynamicRouteConfig
		}

		if dynRouteCfg.Name != l.RouteName {
			continue
		}

		for _, virtualHost := range dynRouteCfg.GetVirtualHosts() {
			for _, domain := range virtualHost.GetDomains() {
				foundAnyRouteDomains = true
				if l.Pod == nil {
					break
				}
				if domain == fmt.Sprintf("%s.%s", l.Pod.Name, l.Pod.Namespace) {
					foundSpecificRouteDomain = true
					break
				}
			}
			if (l.Pod == nil && foundAnyRouteDomains) || foundSpecificRouteDomain {
				break
			}
		}
	}

	if !foundAnyRouteDomains {
		log.Error().Msgf("must have at least one dynamic route config domain: %+v", envoyConfig.Routes.GetDynamicRouteConfigs())
		return ErrNoDynamicRouteConfigDomains
	}

	if l.Pod != nil && !foundSpecificRouteDomain {
		return ErrDynamicRouteConfigDomainNotFound
	}

	return nil
}

// Suggestion implements common.Runnable
func (l RouteDomainCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l RouteDomainCheck) FixIt() error {
	panic("implement me")
}

// Info implements common.Runnable
func (l RouteDomainCheck) Info() string {
	return fmt.Sprintf("Checking whether %s is configured with correct %s Envoy route", l.ConfigGetter.GetObjectName(), l.RouteName)
}

// HasOutboundDynamicRouteConfigDomainCheck creates a new common.Runnable, which checks
// whether the Envoy config has an outbound dynamic route domain to the Pod.
func HasOutboundDynamicRouteConfigDomainCheck(configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		Pod:          pod,
		RouteName:    OutboundDynamicRouteConfigName,
	}
}

// HasInboundDynamicRouteConfigDomainCheck creates a new common.Runnable, which checks
// whether the Envoy config has an inbound dynamic route domain to the Pod.
func HasInboundDynamicRouteConfigDomainCheck(configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		Pod:          pod,
		RouteName:    InboundDynamicRouteConfigName,
	}
}
