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
)
