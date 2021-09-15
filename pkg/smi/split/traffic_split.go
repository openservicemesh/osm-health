package split

import (
	"context"
	"fmt"

	smiSplitClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// Verify interface compliance
var _ runner.Runnable = (*TrafficSplitCheck)(nil)

// TrafficSplitCheck implements common.Runnable
type TrafficSplitCheck struct {
	osmVersion  osm.ControllerVersion
	client      kubernetes.Interface
	pod         *corev1.Pod
	splitClient smiSplitClient.Interface
}

// NewTrafficSplitCheck creates a TrafficSplitCheck which checks whether a pod is affected by an SMI traffic split
func NewTrafficSplitCheck(osmVersion osm.ControllerVersion, client kubernetes.Interface, pod *corev1.Pod, smiSplitClient smiSplitClient.Interface) TrafficSplitCheck {
	return TrafficSplitCheck{
		osmVersion:  osmVersion,
		client:      client,
		pod:         pod,
		splitClient: smiSplitClient,
	}
}

// Description implements common.Runnable
func (check TrafficSplitCheck) Description() string {
	return fmt.Sprintf("Checking whether pod %s participates in a traffic split", check.pod.Name)
}

// Run implements common.Runnable
func (check TrafficSplitCheck) Run() outcomes.Outcome {
	switch osm.SupportedTrafficSplit[check.osmVersion] {
	case osm.V1Alpha2:
		return check.runForTrafficSplitV1alpha2()
	default:
		return outcomes.Fail{Error: fmt.Errorf(
			"OSM Controller version could not be mapped to a TrafficSplit version. Supported versions are v0.5 through v0.9")}
	}
}

func (check TrafficSplitCheck) runForTrafficSplitV1alpha2() outcomes.Outcome {
	ns := check.pod.Namespace
	services, err := pod.GetMatchingServices(check.client, check.pod.ObjectMeta.GetLabels(), ns)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if len(services) == 0 {
		return outcomes.Info{Diagnostics: fmt.Sprintf("pod '%s/%s' does not have a corresponding service", ns, check.pod.Name)}
	}

	trafficSplits, err := check.splitClient.SplitV1alpha2().TrafficSplits(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	for _, trafficSplit := range trafficSplits.Items {
		spec := trafficSplit.Spec
		for _, backend := range spec.Backends {
			for _, svc := range services {
				if backend.Service == svc.Name {
					return outcomes.Info{
						Diagnostics: fmt.Sprintf("pod '%s/%s' participates in traffic split for service '%s/%s'", ns, check.pod.Name, ns, spec.Service),
					}
				}
			}
		}
	}
	return outcomes.Info{
		Diagnostics: fmt.Sprintf("pod '%s/%s' does not participate in any traffic split", check.pod.Namespace, check.pod.Name),
	}
}

// Suggestion implements common.Runnable
func (check TrafficSplitCheck) Suggestion() string {
	return fmt.Sprintf("Check that pod's service corresponds to a TrafficSplit backend. To get TrafficSplits in the namespace, use: \"kubectl get trafficsplit -n %s -o yaml\"", check.pod.Namespace)
}

// FixIt implements common.Runnable
func (check TrafficSplitCheck) FixIt() error {
	panic("implement me")
}
