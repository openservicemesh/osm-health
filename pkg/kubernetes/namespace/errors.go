package namespace

import (
	"errors"
)

var (
	// ErrNotAnnotatedForSidecarInjection is used when an object is expected to have sidecar injection annotation but it does not.
	ErrNotAnnotatedForSidecarInjection = errors.New("not annotated for sidecar injection")

	// ErrNotMonitoredByOSMController is used when namespace is expected to be monitored by OSM but is not.
	ErrNotMonitoredByOSMController = errors.New("not monitored by OSM controller")
)
