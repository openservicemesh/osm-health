package osm

import "errors"

var (
	// ErrorNoControllerPodsExistInNamespace denotes when no osm-controller pods exist in the specified namespace.
	ErrorNoControllerPodsExistInNamespace = errors.New("no osm-controller pods exist in the specified namespace")

	// ErrorControllerNotReady is the error when the osm-controller is not ready,
	// such as when the controller HTTP server readiness probe returns errors.
	ErrorControllerNotReady = errors.New("osm-controller is not ready")

	// ErrorControllerNotAlive is the error when the osm-controller is not alive,
	// such as when the controller HTTP server liveness probe returns errors.
	ErrorControllerNotAlive = errors.New("osm-controller is not alive")
)
