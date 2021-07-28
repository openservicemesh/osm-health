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

	expectedListenerName, exists := osm.ListenerNames[l.ControllerVersion]
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

// Info implements common.Runnable
func (l HasListenerCheck) Info() string {
	return fmt.Sprintf("Checking whether %s is configured with correct Envoy listener", l.ConfigGetter.GetObjectName())
}

// HasListener creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener.
func HasListener(configGetter ConfigGetter, osmVersion osm.ControllerVersion) common.Runnable {
	return HasListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
	}
}
