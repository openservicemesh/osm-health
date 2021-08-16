package podhelper

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/mesh"
)

// Verify interface compliance
var _ common.Runnable = (*ProxyUUIDLabelCheck)(nil)

// ProxyUUIDLabelCheck implements common.Runnable
type ProxyUUIDLabelCheck struct {
	pod *corev1.Pod
}

// HasProxyUUIDLabel checks whether a pod has a valid proxy UUID label which is added when a pod is added to a mesh
func HasProxyUUIDLabel(pod *corev1.Pod) ProxyUUIDLabelCheck {
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
	if !mesh.ProxyLabelExists(*check.pod) {
		return ErrProxyUUIDLabelMissing
	}
	return nil
}

// Suggestion implements common.Runnable
func (check ProxyUUIDLabelCheck) Suggestion() string {
	return "Verify that the pod is in a meshed namespace. Try: \"osm namespace list\""
}

// FixIt implements common.Runnable
func (check ProxyUUIDLabelCheck) FixIt() error {
	panic("implement me")
}
