package osm

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm/pkg/constants"
)

// Verify interface compliance
var _ runner.Runnable = (*NoBadOsmPodLogsCheck)(nil)

// NoBadOsmPodLogsCheck implements common.Runnable
type NoBadOsmPodLogsCheck struct {
	client                   kubernetes.Interface
	osmControlPlaneNamespace common.MeshNamespace
	podLabelSelector         metav1.LabelSelector
	podName                  string
	containerName            string
}

// HasNoBadOsmPodLogsCheck checks whether the specified osm pods in the namespace have bad (fatal/error/warning/fail) log messages
func HasNoBadOsmPodLogsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace, podLabelSelector metav1.LabelSelector, podName string, containerName string) NoBadOsmPodLogsCheck {
	return NoBadOsmPodLogsCheck{
		client:                   client,
		osmControlPlaneNamespace: osmControlPlaneNamespace,
		podLabelSelector:         podLabelSelector,
		podName:                  podName,
		containerName:            containerName,
	}
}

// HasNoBadOsmControllerLogsCheck checks whether the osm controller pods in the namespace have bad (fatal/error/warning/fail) log messages
func HasNoBadOsmControllerLogsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace) NoBadOsmPodLogsCheck {
	return HasNoBadOsmPodLogsCheck(
		client,
		osmControlPlaneNamespace,
		metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMControllerName}},
		constants.OSMControllerName,
		constants.OSMControllerName)
}

// HasNoBadOsmInjectorLogsCheck checks whether the osm controller pods in the namespace have bad (fatal/error/warning/fail) log messages
func HasNoBadOsmInjectorLogsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace) NoBadOsmPodLogsCheck {
	return HasNoBadOsmPodLogsCheck(
		client,
		osmControlPlaneNamespace,
		metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMInjectorName}},
		constants.OSMInjectorName,
		constants.OSMInjectorName)
}

// Description implements common.Runnable
func (check NoBadOsmPodLogsCheck) Description() string {
	return fmt.Sprintf("Checking whether namespace %s has bad (fatal/error/warning/fail) logs in %s pods (container: %s)", check.osmControlPlaneNamespace, check.podName, check.containerName)
}

// Run implements common.Runnable
func (check NoBadOsmPodLogsCheck) Run() outcomes.Outcome {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(check.podLabelSelector.MatchLabels).String(),
	}
	pods, err := check.client.CoreV1().Pods(check.osmControlPlaneNamespace.String()).List(context.TODO(), listOptions)
	if err != nil {
		return outcomes.Fail{Error: fmt.Errorf("unable to list %s pods in namespace %s", check.podName, check.osmControlPlaneNamespace)}
	}

	osmPodErrCount := 0
	for i := range pods.Items {
		if err := podhelper.HasNoBadLogs(check.client, &pods.Items[i], check.podName).GetError(); err != nil {
			osmPodErrCount++
			log.Error().Err(err)
		}
	}

	if osmPodErrCount != 0 {
		return outcomes.Fail{Error: errors.Errorf("%s pods (container: %s) in namespace %s have %d errors", check.podName, check.containerName, check.osmControlPlaneNamespace, osmPodErrCount)}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable.
func (check NoBadOsmPodLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadOsmPodLogsCheck) FixIt() error {
	panic("implement me")
}
