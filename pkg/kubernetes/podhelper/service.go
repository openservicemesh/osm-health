package podhelper

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// Verify interface compliance
var _ runner.Runnable = (*ServiceCheck)(nil)

// ServiceCheck implements common.Runnable
type ServiceCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// NewServiceCheck checks whether a pod has a corresponding service
func NewServiceCheck(client kubernetes.Interface, pod *corev1.Pod) ServiceCheck {
	return ServiceCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (check ServiceCheck) Description() string {
	return fmt.Sprintf("Checking whether destination pod %s has at least one service", check.pod.Name)
}

// Run implements common.Runnable
func (check ServiceCheck) Run() outcomes.Outcome {
	ns := check.pod.Namespace
	services, err := pod.GetMatchingServices(check.client, check.pod.ObjectMeta.GetLabels(), ns)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if len(services) == 0 {
		return outcomes.Fail{Error: errors.Wrapf(ErrNoService, "destination pod '%s/%s' does not have a corresponding service", ns, check.pod.Name)}
	}
	svcNames := []string{}
	for _, s := range services {
		svcNames = append(svcNames, s.ObjectMeta.Name)
	}
	return outcomes.Pass{Msg: fmt.Sprintf("found service(s) %v for destination pod '%s/%s'", svcNames, ns, check.pod.Name)}
}

// Suggestion implements common.Runnable
func (check ServiceCheck) Suggestion() string {
	return "Destination pod should have at least one service associated"
}

// FixIt implements common.Runnable
func (check ServiceCheck) FixIt() error {
	panic("implement me")
}
