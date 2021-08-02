package pod

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/reference"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodNoBadEventsCheck implements common.Runnable
type PodNoBadEventsCheck struct {
	client    kubernetes.Interface
	pod    *v1.Pod
}

// DoesNotHaveBadEvents implements common.Runnable
func DoesNotHaveBadEvents(client kubernetes.Interface, pod *v1.Pod) common.Runnable {
	return PodNoBadEventsCheck{
		client: client,
		pod: pod,
	}
}

// Info implements common.Runnable
func (check PodNoBadEventsCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s has events of type!=Normal", check.pod.Name)
}

// Run implements common.Runnable
func (check PodNoBadEventsCheck) Run() error {
	eventsInterface := check.client.CoreV1().Events(check.pod.Namespace)
	var events *corev1.EventList
	ref, err := reference.GetReference(scheme.Scheme, check.pod)
	if err != nil {
		return fmt.Errorf("unable to construct reference to '%#v': %v", check.pod, err)
	}

	ref.Kind = ""
	if _, isMirrorPod := check.pod.Annotations[corev1.MirrorPodAnnotationKey]; isMirrorPod {
		ref.UID = types.UID(check.pod.Annotations[corev1.MirrorPodAnnotationKey])
	}

	selectorString := "type!=Normal"
	options := metav1.ListOptions{FieldSelector: selectorString}
	if events, err = eventsInterface.List(context.TODO(), options); err != nil {
		return fmt.Errorf("unable to search events of pod '%#v': %v", check.pod, err)
	}

	if len(events.Items) == 0 {
		return nil
	}

	abnormalEvents := "Events:\n  Type\tReason\tTimestamp\tFrom\tMessage\n"
	abnormalEvents += "----\t------\t----\t----\t-------\n"
	for _, event := range events.Items {
		var abnormalEvent string
		source := event.Source.Component
		if source == "" {
			source = event.ReportingController
		}
		abnormalEvent = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", event.Type, event.Reason, event.EventTime, source, strings.TrimSpace(event.Message))
		abnormalEvents += abnormalEvent
	}
	log.Warn().Msg(abnormalEvents)

	return fmt.Errorf("pod '%s' has events that are of type!=Normal - run 'kubectl get events --field-selector %s' to inspect events", check.pod.Name, selectorString)
}
