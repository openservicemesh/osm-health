package osm

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm/pkg/constants"
)

// Verify interface compliance
var _ runner.Runnable = (*NoBadOsmInitLogsCheck)(nil)

// NoBadOsmInitLogsCheck implements common.Runnable
type NoBadOsmInitLogsCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// HasNoBadOsmInitLogsCheck checks whether the osm-init container of the pod has bad (fatal/error/warning/fail) log messages
func HasNoBadOsmInitLogsCheck(client kubernetes.Interface, pod *corev1.Pod) NoBadOsmInitLogsCheck {
	return NoBadOsmInitLogsCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (check NoBadOsmInitLogsCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has bad (fatal/error/warning/fail) logs in osm-init container", check.pod.Name)
}

// Run implements common.Runnable
func (check NoBadOsmInitLogsCheck) Run() outcomes.Outcome {
	return podhelper.HasNoBadLogs(check.client, check.pod, constants.InitContainerName)
}

// Suggestion implements common.Runnable.
func (check NoBadOsmInitLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadOsmInitLogsCheck) FixIt() error {
	panic("implement me")
}
