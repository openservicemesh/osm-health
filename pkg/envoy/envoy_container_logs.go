package envoy

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
)

// Verify interface compliance
var _ common.Runnable = (*NoBadEnvoyLogsCheck)(nil)

// NoBadEnvoyLogsCheck implements common.Runnable
type NoBadEnvoyLogsCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// HasNoBadEnvoyLogsCheck checks whether the envoy container of the pod has bad (fatal/error/warning/fail) log messages
func HasNoBadEnvoyLogsCheck(client kubernetes.Interface, pod *corev1.Pod) NoBadEnvoyLogsCheck {
	return NoBadEnvoyLogsCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (check NoBadEnvoyLogsCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has bad (fatal/error/warning/fail) logs in envoy container", check.pod.Name)
}

// Run implements common.Runnable
func (check NoBadEnvoyLogsCheck) Run() outcomes.Outcome {
	return podhelper.HasNoBadLogs(check.client, check.pod, "envoy")
}

// Suggestion implements common.Runnable.
func (check NoBadEnvoyLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadEnvoyLogsCheck) FixIt() error {
	panic("implement me")
}
