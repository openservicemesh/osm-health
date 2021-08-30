package envoy

import (
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// Verify interface compliance
var _ common.Runnable = (*RouteDomainCheck)(nil)

// RouteDomainCheck implements common.Runnable
type RouteDomainCheck struct {
	*corev1.Pod
	ConfigGetter
	RouteName string
	Domain    string
	Client    kubernetes.Interface
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

	var possibleDomains []string
	if len(l.Domain) > 0 {
		possibleDomains = append(possibleDomains, l.Domain)
	}
	if l.Client != nil {
		svcs, err := kuberneteshelper.GetMatchingServices(l.Client, l.Labels, l.Namespace)
		if err != nil {
			return outcomes.FailedOutcome{Error: err}
		}
		for _, svc := range svcs {
			possibleDomains = append(possibleDomains, svc.Name+"."+svc.Namespace)
		}
	}

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
				for _, possibleDomain := range possibleDomains {
					if domain == possibleDomain {
						foundSpecificRouteDomain = true
						break
					}
				}
				if foundSpecificRouteDomain {
					break
				}
			}
			if (len(possibleDomains) == 0 && foundAnyRouteDomains) || foundSpecificRouteDomain {
				break
			}
		}
	}

	if !foundAnyRouteDomains {
		log.Error().Msgf("must have at least one dynamic route config domain: %+v", envoyConfig.Routes.GetDynamicRouteConfigs())
		return outcomes.FailedOutcome{Error: ErrNoDynamicRouteConfigDomains}
	}

	if len(possibleDomains) > 0 && !foundSpecificRouteDomain {
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
func NewOutboundRouteDomainPodCheck(client kubernetes.Interface, configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		Pod:          pod,
		ConfigGetter: configGetter,
		RouteName:    OutboundDynamicRouteConfigName,
		Client:       client,
	}
}

// NewInboundRouteDomainPodCheck creates a DestinationEndpointCheck which checks whether the Envoy config has an inbound dynamic route domain to the Pod
func NewInboundRouteDomainPodCheck(client kubernetes.Interface, configGetter ConfigGetter, pod *corev1.Pod) RouteDomainCheck {
	return RouteDomainCheck{
		Pod:          pod,
		ConfigGetter: configGetter,
		RouteName:    InboundDynamicRouteConfigName,
		Client:       client,
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
