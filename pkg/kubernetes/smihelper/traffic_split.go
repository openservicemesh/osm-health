package smihelper

import (
	"context"
	"fmt"

	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

// Verify interface compliance
var _ common.Runnable = (*TrafficSplitCheck)(nil)

// TrafficSplitCheck implements common.Runnable
type TrafficSplitCheck struct {
	client          kubernetes.Interface
	pod             *v1.Pod
	smiAccessClient smiAccessClient.Interface
}

// IsInTrafficSplit checks whether a pod is affected by an SMI traffic split
func IsInTrafficSplit(client kubernetes.Interface, pod *v1.Pod, smiAccessClient smiAccessClient.Interface) TrafficSplitCheck {
	return TrafficSplitCheck{
		client:          client,
		pod:             pod,
		smiAccessClient: smiAccessClient,
	}
}

// Info implements common.Runnable
func (check TrafficSplitCheck) Info() string {
	return fmt.Sprintf("Checking whether pod %s participates in a traffic split", check.pod.Name)
}

// Run implements common.Runnable
func (check TrafficSplitCheck) Run() error {
	services, err := kuberneteshelper.GetMatchingServices(check.client, check.pod)
	if err != nil {
		return ErrGettingServices
	}
	trafficSplits, err := check.smiAccessClient.SplitV1alpha2().TrafficSplits(check.pod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return ErrGettingTrafficSplits
	}
	for _, trafficSplit := range trafficSplits.Items {
		spec := trafficSplit.Spec
		for _, backend := range spec.Backends {
			for _, svc := range services {
				if backend.Service == svc.Name {
					return nil
				}
			}
		}
	}
	return ErrNoTrafficSplitForPod
}

// Suggestion implements common.Runnable
func (check TrafficSplitCheck) Suggestion() string {
	return fmt.Sprintf("Check that pod's service corresponds to a TrafficSplit backend. To get TrafficSplits in the namespace, use: \"kubectl get trafficsplit -n %s -o yaml\"", check.pod.Namespace)
}

// FixIt implements common.Runnable
func (check TrafficSplitCheck) FixIt() error {
	panic("implement me")
}
