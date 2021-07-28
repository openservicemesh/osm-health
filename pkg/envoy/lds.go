package envoy

import (
	"errors"

	"github.com/openservicemesh/osm-health/pkg/common"
	k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"
)

// HasListener implements common.Runnable
type HasListener struct {
	k8s.Pod
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
func NewHasListener(pod k8s.Pod) common.Runnable {
	return HasListener{
		Pod: pod,
	}
}
