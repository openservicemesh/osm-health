package pod

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// Verify interface compliance
var _ common.Runnable = (*NoBadEventsCheck)(nil)

// NoBadEventsCheck implements common.Runnable
type NoBadEventsCheck struct {
	client kubernetes.Interface
	pod    *v1.Pod
}

// DoesNotHaveBadEvents checks whether a pod has abnormal (type!=Normal) events
func DoesNotHaveBadEvents(client kubernetes.Interface, pod *v1.Pod) NoBadEventsCheck {
	return NoBadEventsCheck{
		client: client,
		pod:    pod,
	}
}

// Info implements common.Runnable
func (check NoBadEventsCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has events of type!=Normal", check.pod.Name)
}

// Run implements common.Runnable
func (check NoBadEventsCheck) Run() error {
	eventsInterface := check.client.CoreV1().Events(check.pod.Namespace)
	var events *corev1.EventList

	selectorString := "type!=Normal"
	options := metav1.ListOptions{FieldSelector: selectorString}
	events, err := eventsInterface.List(context.TODO(), options)
	if err != nil {
		return fmt.Errorf("unable to search events of pod '%#v': %v", check.pod, err)
	}

	if len(events.Items) == 0 {
		return nil
	}

	return fmt.Errorf("pod '%s' has events that are of 'type!=Normal' - run 'kubectl get events --field-selector %s' to inspect events", check.pod.Name, selectorString)
}

// Suggestion implements common.Runnable.
func (check NoBadEventsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadEventsCheck) FixIt() error {
	panic("implement me")
}
