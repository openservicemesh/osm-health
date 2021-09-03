package envoy

import (
	"fmt"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// Verify interface compliance
var _ common.Runnable = (*ListenerCheck)(nil)

// ListenerCheck implements common.Runnable
type ListenerCheck struct {
	ConfigGetter
	osm.ControllerVersion

	// This is used for Info() function from Runnable. Helps the logs identify what kind of a listener we are looking for.
	listenerType string

	// This map will be used to get the expected NAME of the Envoy listener depending on the OSM version in use.
	expectedListenersPerVersion map[osm.ControllerVersion]string
}

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
func NewOutboundListenerCheck(configGetter ConfigGetter, osmVersion osm.ControllerVersion) ListenerCheck {
	return ListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "outbound",

		expectedListenersPerVersion: osm.OutboundListenerNames,
	}
}

// NewInboundListenerCheck creates a ListenerCheck which checks whether the given Pod has an Envoy with properly configured listener.
func NewInboundListenerCheck(configGetter ConfigGetter, osmVersion osm.ControllerVersion) ListenerCheck {
	return ListenerCheck{
		ConfigGetter:      configGetter,
		ControllerVersion: osmVersion,
		listenerType:      "inbound",

		expectedListenersPerVersion: osm.InboundListenerNames,
	}
}
