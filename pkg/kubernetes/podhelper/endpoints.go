package podhelper

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

// EndpointsCheck implements common.Runnable
type EndpointsCheck struct {
	client kubernetes.Interface
	pod    *corev1.Pod
}

// NewEndpointsCheck creates an EndpointsCheck which checks whether a pod is
// referenced by a Kubernetes Endpoints resource.
func NewEndpointsCheck(client kubernetes.Interface, pod *corev1.Pod) EndpointsCheck {
	return EndpointsCheck{
		client: client,
		pod:    pod,
	}
}

// Description implements common.Runnable
func (e EndpointsCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s/%s is referenced by a Kubernetes Endpoints resource", e.pod.Namespace, e.pod.Name)
}

// Run implements common.Runnable
func (e EndpointsCheck) Run() outcomes.Outcome {
	eps, err := e.client.CoreV1().Endpoints(e.pod.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	found := false
	for _, ep := range eps.Items {
		for _, subset := range ep.Subsets {
			for _, addr := range subset.Addresses {
				if addr.TargetRef.Name == e.pod.Name {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return outcomes.Fail{Error: ErrPodNotInEndpoints}
	}
	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (e EndpointsCheck) Suggestion() string {
	return fmt.Sprintf("Verify the selector on the Kubernetes Service in the %s namespace (kubectl get svc -n %s <name> -o jsonpath='{.spec.selector}') that should be backed by Pod %s matches the Pod's labels (kubectl get pod -n %s %s -o jsonpath='{.metadata.labels}').", e.pod.Namespace, e.pod.Namespace, e.pod.Name, e.pod.Namespace, e.pod.Name)
}

// FixIt implements common.Runnable
func (e EndpointsCheck) FixIt() error {
	panic("implement me")
}
