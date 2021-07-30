package envoy

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*HasDestinationEndpointsCheck)(nil)

// HasDestinationEndpointsCheck implements common.Runnable
type HasDestinationEndpointsCheck struct {
	*v1.Pod
	ConfigGetter
}

// Run implements common.Runnable
func (l HasDestinationEndpointsCheck) Run() error {
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

	// If Pod was defined -- check if this pod IP is in the list of endpoints.
	if l.Pod != nil {
		foundIt := false
		// Check for a specific Pod.
		for _, ept := range cla.Endpoints[0].LbEndpoints {
			if ept.GetEndpoint().Address.GetSocketAddress().Address == l.Status.PodIP {
				foundIt = true
				break
			}
		}
		if !foundIt {
			return ErrEndpointNotFound
		}
	}

	return nil
}

// Info implements common.Runnable
func (l HasDestinationEndpointsCheck) Info() string {
	txt := "at least one destination"
	if l.Pod != nil {
		txt = fmt.Sprintf("%s as a destination", l.Status.PodIP)
	}

	return fmt.Sprintf("Checking whether %s is configured with %s endpoint", l.ConfigGetter.GetObjectName(), txt)
}

// HasDestinationEndpoints creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener for the local payload.
func HasDestinationEndpoints(configGetter ConfigGetter) common.Runnable {
	return HasDestinationEndpointsCheck{
		ConfigGetter: configGetter,
	}
}

// HasSpecificEndpoint creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener for the local payload.
func HasSpecificEndpoint(configGetter ConfigGetter, pod *v1.Pod) common.Runnable {
	return HasDestinationEndpointsCheck{
		ConfigGetter: configGetter,
		Pod:          pod,
	}
}
