package envoy

import "errors"

var (
	// ErrEnvoyListenerMissing is an error
	ErrEnvoyListenerMissing = errors.New("envoy listener missing")

	// ErrEnvoyConfigEmpty is an error
	ErrEnvoyConfigEmpty = errors.New("envoy config is empty")

	// ErrOSMControllerVersionUnrecognized is an error
	ErrOSMControllerVersionUnrecognized = errors.New("osm controller version not recognized")

	// ErrIncorrectlyInitializedConfigGetter is an error
	ErrIncorrectlyInitializedConfigGetter = errors.New("incorrectly initialized config getter")

	// ErrNoDestinationEndpoints is an error
	ErrNoDestinationEndpoints = errors.New("no destination endpoints")

	// ErrUnmarshalingClusterLoadAssigment is an error
	ErrUnmarshalingClusterLoadAssigment = errors.New("error unmarshaling envoy cluster load assigment")

	// ErrEndpointNotFound is an error
	ErrEndpointNotFound = errors.New("endpoint not found")
)
