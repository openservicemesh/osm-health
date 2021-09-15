package envoy

import (
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	pod "github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// Verify interface compliance
var _ runner.Runnable = (*RouteDomainCheck)(nil)

// RouteDomainCheck implements common.Runnable
type RouteDomainCheck struct {
	ConfigGetter
	RouteName string
	Domains   map[string]bool
}

// Run implements common.Runnable
func (check RouteDomainCheck) Run() outcomes.Outcome {
	if check.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}

	envoyConfig, err := check.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	foundAnyRouteDomains := false
	foundSpecificRouteDomain := false

	for _, rawDynRouteCfg := range envoyConfig.Routes.GetDynamicRouteConfigs() {
		var dynRouteCfg envoy_config_route_v3.RouteConfiguration
		if err = rawDynRouteCfg.GetRouteConfig().UnmarshalTo(&dynRouteCfg); err != nil {
			return outcomes.Fail{Error: ErrUnmarshalingDynamicRouteConfig}
		}

		if dynRouteCfg.Name != check.RouteName {
			continue
		}

		for _, virtualHost := range dynRouteCfg.GetVirtualHosts() {
			for _, domain := range virtualHost.GetDomains() {
				foundAnyRouteDomains = true
				if len(check.Domains) == 0 {
					break
				}
				if _, ok := check.Domains[domain]; ok {
					foundSpecificRouteDomain = true
					break
				}
			}
			if (len(check.Domains) == 0 && foundAnyRouteDomains) || foundSpecificRouteDomain {
				break
			}
		}
	}

	if !foundAnyRouteDomains {
		log.Error().Msgf("must have at least one dynamic route config domain in envoy config")
		return outcomes.Fail{Error: ErrNoDynamicRouteConfigDomains}
	}

	if len(check.Domains) > 0 && !foundSpecificRouteDomain {
		return outcomes.Fail{Error: ErrDynamicRouteConfigDomainNotFound}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (check RouteDomainCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (check RouteDomainCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (check RouteDomainCheck) Description() string {
	return fmt.Sprintf("Checking whether %s is configured with correct %s Envoy route", check.ConfigGetter.GetObjectName(), check.RouteName)
}

// NewOutboundRouteDomainPodCheck creates a new common.Runnable, which checks
// whether the Envoy config has outbound dynamic route domains to the Pod's services.
func NewOutboundRouteDomainPodCheck(client kubernetes.Interface, configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return NewPodServicesRouteDomainCheck(client, configGetter, pod, OutboundDynamicRouteConfigName)
}

// NewInboundRouteDomainPodCheck creates a new common.Runnable, which checks
// whether the Envoy config has inbound dynamic route domains from the Pod's services.
func NewInboundRouteDomainPodCheck(client kubernetes.Interface, configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return NewPodServicesRouteDomainCheck(client, configGetter, pod, InboundDynamicRouteConfigName)
}

// NewPodServicesRouteDomainCheck checks whether the pod's corresponding service's domains are
// contained in the envoy dynamic route config domain list.
func NewPodServicesRouteDomainCheck(client kubernetes.Interface, configGetter ConfigGetter, podToCheck *corev1.Pod, routeName string) RouteDomainCheck {
	podSvcs, err := pod.GetMatchingServices(client, podToCheck.ObjectMeta.GetLabels(), podToCheck.Namespace)
	if err != nil {
		log.Warn().Msgf("unable to obtain the services of pod %s/%s", podToCheck.Namespace, podToCheck.Name)
	}

	domains := make(map[string]bool)
	for _, svc := range podSvcs {
		domains[fmt.Sprintf("%s.%s", svc.Name, svc.Namespace)] = false
	}

	return RouteDomainCheck{
		ConfigGetter: configGetter,
		RouteName:    routeName,
		Domains:      domains,
	}
}

// NewOutboundRouteDomainHostCheck creates a DestinationEndpointCheck which checks whether the Envoy config has an outbound dynamic route domain to the URL
func NewOutboundRouteDomainHostCheck(configGetter ConfigGetter, destinationHost string) RouteDomainCheck {
	return RouteDomainCheck{
		ConfigGetter: configGetter,
		RouteName:    OutboundDynamicRouteConfigName,
		Domains:      map[string]bool{destinationHost: true},
	}
}
