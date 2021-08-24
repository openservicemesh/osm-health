package envoy

import (
	"fmt"

	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Verify interface compliance
var _ common.Runnable = (*DestinationEndpointChecker)(nil)

// DestinationEndpointChecker implements common.Runnable
type DestinationEndpointChecker struct {
	*corev1.Pod
	ConfigGetter
}

// Run implements common.Runnable
func (l DestinationEndpointChecker) Run() outcomes.Outcome {
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

	foundAnyEndpoints := false
	// If Pod was defined -- check if this pod IP is in the list of endpoints.
	foundSpecificEndpoint := false

	for _, dynEpt := range envoyConfig.Endpoints.GetDynamicEndpointConfigs() {
		var cla envoy_config_endpoint_v3.ClusterLoadAssignment
		if err = dynEpt.GetEndpointConfig().UnmarshalTo(&cla); err != nil {
			return outcomes.FailedOutcome{Error: ErrUnmarshalingClusterLoadAssigment}
		}

		for _, ept := range cla.GetEndpoints() {
			for _, lbEpt := range ept.GetLbEndpoints() {
				foundAnyEndpoints = true
				if l.Pod == nil {
					break
				}
				if lbEpt.GetEndpoint().GetAddress().GetSocketAddress().GetAddress() == l.Status.PodIP {
					foundSpecificEndpoint = true
					break
				}
			}
			if (l.Pod == nil && foundAnyEndpoints) || foundSpecificEndpoint {
				break
			}
		}
	}

	if !foundAnyEndpoints {
		log.Error().Msgf("must have at least one destination endpoint: %+v", envoyConfig.Endpoints.GetDynamicEndpointConfigs())
		return outcomes.FailedOutcome{Error: ErrNoDestinationEndpoints}
	}

	if l.Pod != nil && !foundSpecificEndpoint {
		return outcomes.FailedOutcome{Error: ErrEndpointNotFound}
	}

	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable
func (l DestinationEndpointChecker) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l DestinationEndpointChecker) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (l DestinationEndpointChecker) Description() string {
	txt := "at least one destination"
	if l.Pod != nil {
		txt = fmt.Sprintf("%s as a destination", l.Status.PodIP)
	}

	return fmt.Sprintf("Checking whether %s is configured with %s endpoint", l.ConfigGetter.GetObjectName(), txt)
}

// HasDestinationEndpoints creates a new common.Runnable, which checks whether
// the given Pod has an Envoy with any endpoints configured.
func HasDestinationEndpoints(configGetter ConfigGetter) DestinationEndpointChecker {
	return DestinationEndpointChecker{
		ConfigGetter: configGetter,
	}
}

// HasSpecificEndpoint creates a new common.Runnable, which checks whether the
// given Pod has an Envoy with an endpoint configured mapping to a specific
// destination Pod.
func HasSpecificEndpoint(configGetter ConfigGetter, pod *corev1.Pod) DestinationEndpointChecker {
	return DestinationEndpointChecker{
		ConfigGetter: configGetter,
		Pod:          pod,
	}
}
