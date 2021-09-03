package podhelper

import (
	"errors"
)

var (
	// ErrExpectedEnvoyImageMissing is used when a pod is expected to have a container with an envoy sidecar image but does not
	ErrExpectedEnvoyImageMissing = errors.New("expected envoy container image missing")

	// ErrExpectedOsmInitImageMissing is used when a pod is expected to have an init container with an osm init image but does not
	ErrExpectedOsmInitImageMissing = errors.New("expected osm init container image missing")

	// ErrExpectedMinNumContainers is used when a pod is expected to have a container with an envoy sidecar image but does not
	ErrExpectedMinNumContainers = errors.New("fewer containers than expected in pod")

	// ErrProxyUUIDLabelMissing is used when a pod is expected to have a valid proxy UUID label but does not
	ErrProxyUUIDLabelMissing = errors.New("pod does not have expected valid proxy UUID label")

	// ErrPodDoesNotHaveContainer is used when a pod does not have a container in the pod spec container list.
	ErrPodDoesNotHaveContainer = errors.New("pod does not have container in pod spec container list")

	// ErrPodNotInEndpoints is used when a pod is expected to be referenced by any Kubernetes Endpoints resources but is not
	ErrPodNotInEndpoints = errors.New("pod not referenced by any Kubernetes Endpoints resources")
)
