package pod

import (
	"fmt"

	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

// Verify interface compliance
var _ common.Runnable = (*ProxyUUIDLabelCheck)(nil)

// ProxyUUIDLabelCheck implements common.Runnable
type ProxyUUIDLabelCheck struct {
	pod *v1.Pod
}

// HasProxyUUIDLabel checks whether a pod has a valid proxy UUID label which is added when a pod is added to a mesh
func HasProxyUUIDLabel(pod *v1.Pod) ProxyUUIDLabelCheck {
	return ProxyUUIDLabelCheck{
		pod: pod,
	}
}

// Info implements common.Runnable
func (check ProxyUUIDLabelCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has a valid proxy UUID label", check.pod.Name)
}

// Run implements common.Runnable
func (check ProxyUUIDLabelCheck) Run() error {
	if !proxyLabelExists(check.pod) {
		return ErrProxyUUIDLabelMissing
	}
	return nil
}

// Suggestion implements common.Runnable
func (check ProxyUUIDLabelCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (check ProxyUUIDLabelCheck) FixIt() error {
	panic("implement me")
}

// TODO: replace with function from osm pkg once it's made public
// proxyLabelExists returns a boolean indicating whether the pod has a proxy UUID label. A proxy UUID label is added when a pod is added to a mesh
func proxyLabelExists(pod *v1.Pod) bool {
	// osm-controller adds a unique label to each pod when it is added to a mesh
	proxyUUID, proxyLabelSet := pod.Labels[constants.EnvoyUniqueIDLabelName]
	return proxyLabelSet && isValidUUID(proxyUUID)
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
