package envoy

import "errors"

var (
	// ErrEnvoyListenerMissing is an error returned when an Envoy does not have a required listener.
	ErrEnvoyListenerMissing = errors.New("envoy listener missing")

	// ErrEnvoyConfigEmpty is an error returned when an Envoy config is completely missing.
	ErrEnvoyConfigEmpty = errors.New("envoy config is empty")

	// ErrOSMControllerVersionUnrecognized is an error returned when the supplied OSM Controller version is not recognized.
	ErrOSMControllerVersionUnrecognized = errors.New("osm controller version not recognized")

	// ErrIncorrectlyInitializedConfigGetter is an error returned when the ConfigGetter struct is not correctly initialized.
	ErrIncorrectlyInitializedConfigGetter = errors.New("incorrectly initialized config getter")

	// ErrNoDestinationEndpoints is an error returned when an Envoy has no destination endpoints.
	ErrNoDestinationEndpoints = errors.New("no destination endpoints")

	// ErrUnmarshalingClusterLoadAssigment is an error returned when the unmarshaling of the Envoy ClusterLoadAssignment struct fails.
	ErrUnmarshalingClusterLoadAssigment = errors.New("error unmarshaling envoy cluster load assigment")

	// ErrEndpointNotFound is an error returned when a specific endpoint is not found in Envoy EDS config.
	ErrEndpointNotFound = errors.New("endpoint not found")
)
