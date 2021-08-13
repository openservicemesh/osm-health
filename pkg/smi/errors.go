package smi

import (
	"errors"
)

var (
	// ErrNoTrafficSplitForPod is used when a pod does not participate in a traffic split
	ErrNoTrafficSplitForPod = errors.New("pod does not participate in a traffic split")
)
