package podhelper

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm/pkg/mesh"
)

// Verify interface compliance
var _ common.Runnable = (*ProxyUUIDLabelCheck)(nil)

// ProxyUUIDLabelCheck implements common.Runnable
type ProxyUUIDLabelCheck struct {
	pod *corev1.Pod
}

// NewProxyUUIDLabelCheck creates a ProxyUUIDLabelCheck which checks whether a pod has a valid proxy UUID label which is added when a pod is added to a mesh
func NewProxyUUIDLabelCheck(pod *corev1.Pod) ProxyUUIDLabelCheck {
	return ProxyUUIDLabelCheck{
		pod: pod,
	}
}

// Description implements common.Runnable
func (check ProxyUUIDLabelCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has a valid proxy UUID label", check.pod.Name)
}

// Run implements common.Runnable
func (check ProxyUUIDLabelCheck) Run() outcomes.Outcome {
	if !mesh.ProxyLabelExists(*check.pod) {
		return outcomes.FailedOutcome{Error: ErrProxyUUIDLabelMissing}
	}
	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable
func (check ProxyUUIDLabelCheck) Suggestion() string {
	return "Verify that the pod is in a meshed namespace. Try: \"osm namespace list\""
}

// FixIt implements common.Runnable
func (check ProxyUUIDLabelCheck) FixIt() error {
	panic("implement me")
}
