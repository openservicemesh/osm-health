package envoy

import (
	"fmt"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// Verify interface compliance
var _ common.Runnable = (*HasListenerCheck)(nil)

// HasListenerCheck implements common.Runnable
type HasListenerCheck struct {
	ConfigGetter
	osm.ControllerVersion

	// This is used for Info() function from Runnable. Helps the logs identify what kind of a listener we are looking for.
	listenerType string

	// This map will be used to get the expected NAME of the Envoy listener depending on the OSM version in use.
	expectedListenersPerVersion map[osm.ControllerVersion]string
}

// Run implements common.Runnable
func (l HasListenerCheck) Run() error {
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

	expectedListenerName, exists := l.expectedListenersPerVersion[l.ControllerVersion]
	if !exists {
		return ErrOSMControllerVersionUnrecognized
	}

	if envoyConfig == nil || envoyConfig.Listeners.DynamicListeners == nil {
		return ErrEnvoyConfigEmpty
	}

	actualListener := envoyConfig.Listeners.DynamicListeners[0]

	if expectedListenerName != actualListener.Name {
		log.Error().Msgf("must have listener with name %s but it is instead %s", expectedListenerName, actualListener.Name)
		return ErrEnvoyListenerMissing
	}

	return nil
}

// Suggestion implements common.Runnable
func (l HasListenerCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l HasListenerCheck) FixIt() error {
	panic("implement me")
}

// Info implements common.Runnable
func (l HasListenerCheck) Info() string {
	return fmt.Sprintf("Checking whether %s is configured with correct %s Envoy listener", l.ConfigGetter.GetObjectName(), l.listenerType)
}

// HasOutboundListener creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener.
func HasOutboundListener(configGetter ConfigGetter, osmVersion osm.ControllerVersion) HasListenerCheck {
	return HasListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "outbound",

		expectedListenersPerVersion: osm.OutboundListenerNames,
	}
}

// HasInboundListener creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener.
func HasInboundListener(configGetter ConfigGetter, osmVersion osm.ControllerVersion) HasListenerCheck {
	return HasListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "inbound",

		expectedListenersPerVersion: osm.InboundListenerNames,
	}
}
