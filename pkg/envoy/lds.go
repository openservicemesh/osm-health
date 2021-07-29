package envoy

import (
	"errors"

	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*HasListener)(nil)

// HasListener implements common.Runnable
type HasListener struct {
	*v1.Pod
}

// Run implements common.Runnable
func (l HasListener) Run() error {
	envoyConfig, err := GetEnvoyConfig(l.Pod)
	if err != nil {
		return err
	}
	if envoyConfig == nil {
		return errors.New("envoy config is empty")
	}
	return nil
}

// Info implements common.Runnable
func (l HasListener) Info() string {
	return ""
}

// NewHasListener creates a new common.Runnable, which checks whether the given Pod has an Envoy with properly configured listener.
func NewHasListener(pod *v1.Pod) HasListener {
	return HasListener{
		Pod: pod,
	}
}
