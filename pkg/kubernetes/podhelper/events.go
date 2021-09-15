package podhelper

import (
	"context"
	"fmt"

	"github.com/openservicemesh/osm-health/pkg/runner"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// Verify interface compliance
var _ runner.Runnable = (*PodEventsCheck)(nil)

// PodEventsCheck implements common.Runnable
type PodEventsCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// NewPodEventsCheck creates a PodEventsCheck which checks whether a pod has abnormal (type!=Normal) events
func NewPodEventsCheck(client kubernetes.Interface, pod *corev1.Pod) PodEventsCheck {
	return PodEventsCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (check PodEventsCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s has events of type!=Normal", check.pod.Name)
}

// Run implements common.Runnable
func (check PodEventsCheck) Run() outcomes.Outcome {
	eventsInterface := check.client.CoreV1().Events(check.pod.Namespace)
	var events *corev1.EventList

	selectorString := "type!=Normal,involvedObject.apiVersion=v1,involvedObject.kind=Pod,involvedObject.name=" + check.pod.Name
	options := metav1.ListOptions{FieldSelector: selectorString}
	events, err := eventsInterface.List(context.TODO(), options)
	if err != nil {
		return outcomes.Fail{Error: fmt.Errorf("unable to search events of pod '%#v': %v", check.pod, err)}
	}

	if len(events.Items) == 0 {
		return outcomes.Pass{}
	}

	return outcomes.Fail{Error: fmt.Errorf("pod '%s' has events that are of 'type!=Normal' - run 'kubectl get events --namespace %s --field-selector %s' to inspect events", check.pod.Name, check.pod.Namespace, selectorString)}
}

// Suggestion implements common.Runnable.
func (check PodEventsCheck) Suggestion() string {
	return fmt.Sprintf("To inspect for unexpected events, try \"kubectl get events --namespace %s --field-selector type!=Normal\"", check.pod.Namespace)
}

// FixIt implements common.Runnable.
func (check PodEventsCheck) FixIt() error {
	panic("implement me")
}
