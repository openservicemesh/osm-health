package envoy

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*DestinationEndpointChecker)(nil)

// DestinationEndpointChecker implements common.Runnable
type DestinationEndpointChecker struct {
	*v1.Pod
	ConfigGetter
}

// Run implements common.Runnable
func (l DestinationEndpointChecker) Run() error {
	if l.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return ErrIncorrectlyInitializedConfigGetter
	}
	envoyConfig, err := l.ConfigGetter.GetConfig()
	if err != nil {
		return err
	}

	if envoyConfig == nil || envoyConfig.Listeners.DynamicListeners == nil {
		return ErrEnvoyConfigEmpty
	}

	if len(envoyConfig.Endpoints.GetDynamicEndpointConfigs()) <= 0 {
		log.Error().Msgf("must have at least one destination endpoint: %+v", envoyConfig.Endpoints.GetDynamicEndpointConfigs())
		return ErrNoDestinationEndpoints
	}

	var cla envoy_config_endpoint_v3.ClusterLoadAssignment
	if err = envoyConfig.Endpoints.GetDynamicEndpointConfigs()[0].GetEndpointConfig().UnmarshalTo(&cla); err != nil {
		return ErrUnmarshalingClusterLoadAssigment
	}

	if len(cla.Endpoints) <= 0 || len(cla.Endpoints[0].LbEndpoints) <= 0 {
		log.Error().Msg("must have at least one destination endpoint")
		return ErrNoDestinationEndpoints
	}

	if l.Pod == nil {
		return nil
	}

	// If Pod was defined -- check if this pod IP is in the list of endpoints.
	foundIt := false
	// Check for a specific Pod.
	for _, ept := range cla.GetEndpoints() {
		for _, lbEpt := range ept.GetLbEndpoints() {
			if lbEpt.GetEndpoint().GetAddress().GetSocketAddress().GetAddress() == l.Status.PodIP {
				foundIt = true
				break
			}
		}
		if foundIt {
			break
		}
	}
	if !foundIt {
		return ErrEndpointNotFound
	}

	return nil
}

// Suggestion implements common.Runnable
func (l DestinationEndpointChecker) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l DestinationEndpointChecker) FixIt() error {
	panic("implement me")
}

// Info implements common.Runnable
func (l DestinationEndpointChecker) Info() string {
	txt := "at least one destination"
	if l.Pod != nil {
		txt = fmt.Sprintf("%s as a destination", l.Status.PodIP)
	}

	return fmt.Sprintf("Checking whether %s is configured with %s endpoint", l.ConfigGetter.GetObjectName(), txt)
}

// HasDestinationEndpoints creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener for the local payload.
func HasDestinationEndpoints(configGetter ConfigGetter) DestinationEndpointChecker {
	return DestinationEndpointChecker{
		ConfigGetter: configGetter,
	}
}

// HasSpecificEndpoint creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener for the local payload.
func HasSpecificEndpoint(configGetter ConfigGetter, pod *v1.Pod) DestinationEndpointChecker {
	return DestinationEndpointChecker{
		ConfigGetter: configGetter,
		Pod:          pod,
	}
}
