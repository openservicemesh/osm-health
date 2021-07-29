package pod

import (
	"errors"
)

var (
	// ErrExpectedEnvoySidcarMissing is used when a pod is expected to have a container with an envoy sidecar image but does not
	ErrExpectedEnvoySidcarMissing = errors.New("expected envoy container image missing")
)
