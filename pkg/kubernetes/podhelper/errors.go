package podhelper

import (
	"errors"
)

var (
	// ErrExpectedEnvoyImageMissing is used when a pod is expected to have a container with an envoy sidecar image but does not
	ErrExpectedEnvoyImageMissing = errors.New("expected envoy container image missing")
)

var (
	// ErrExpectedMinNumContainers is used when a pod is expected to have a container with an envoy sidecar image but does not
	ErrExpectedMinNumContainers = errors.New("fewer containers than expected in pod")
)

var (
	// ErrProxyUUIDLabelMissing is used when a pod is expected to have a valid proxy UUID label but does not
	ErrProxyUUIDLabelMissing = errors.New("pod does not have expected valid proxy UUID label")
)
