package smihelper

import (
	"errors"
)

var (
	// ErrGettingServices is used when a services for a given pod's namespace cannot be found
	ErrGettingServices = errors.New("services for pod's namespace could not be found")

	// ErrGettingTrafficSplits is used when a traffic splits cannot be listed for a given namespace
	ErrGettingTrafficSplits = errors.New("traffic splits could not be found")

	// ErrNoTrafficSplitForPod is used when a pod does not participate in a traffic split
	ErrNoTrafficSplitForPod = errors.New("pod does not participate in a traffic split")
)
