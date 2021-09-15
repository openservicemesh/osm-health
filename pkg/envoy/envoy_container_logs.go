package envoy

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// Verify interface compliance
var _ runner.Runnable = (*BadLogsCheck)(nil)

// BadLogsCheck implements common.Runnable
type BadLogsCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// NewBadLogsCheck creates an BadLogsCheck which checks whether the envoy container of the pod has bad (fatal/error/warning/fail) log messages
func NewBadLogsCheck(client kubernetes.Interface, pod *corev1.Pod) BadLogsCheck {
	return BadLogsCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (check BadLogsCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has bad (fatal/error/warning/fail) logs in envoy container", check.pod.Name)
}

// Run implements common.Runnable
func (check BadLogsCheck) Run() outcomes.Outcome {
	return podhelper.HasNoBadLogs(check.client, check.pod, "envoy")
}

// Suggestion implements common.Runnable.
func (check BadLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check BadLogsCheck) FixIt() error {
	panic("implement me")
}
