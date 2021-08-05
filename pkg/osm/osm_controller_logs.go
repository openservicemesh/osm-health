package osm

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/podhelper"
	"github.com/openservicemesh/osm/pkg/constants"
)

// Verify interface compliance
var _ common.Runnable = (*NoBadOsmControllerLogsCheck)(nil)

// NoBadOsmControllerLogsCheck implements common.Runnable
type NoBadOsmControllerLogsCheck struct {
	client                   kubernetes.Interface
	osmControlPlaneNamespace string
}

// HasNoBadOsmControllerLogsCheck checks whether the osm controller pods in the controller namespace have bad (fatal/error/warning/fail) log messages
func HasNoBadOsmControllerLogsCheck(client kubernetes.Interface, osmControlPlaneNamespace string) NoBadOsmControllerLogsCheck {
	return NoBadOsmControllerLogsCheck{
		client:                   client,
		osmControlPlaneNamespace: osmControlPlaneNamespace,
	}
}

// Info implements common.Runnable
func (check NoBadOsmControllerLogsCheck) Info() string {
	return fmt.Sprintf("Checking whether namespace %s has bad (fatal/error/warning/fail) logs in osm controller pods", check.osmControlPlaneNamespace)
}

// Run implements common.Runnable
func (check NoBadOsmControllerLogsCheck) Run() error {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMControllerName}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	pods, err := check.client.CoreV1().Pods(check.osmControlPlaneNamespace).List(context.TODO(), listOptions)
	if err != nil {
		return fmt.Errorf("unable to list osm controller pods in namespace %s", check.osmControlPlaneNamespace)
	}

	for i := range pods.Items {
		if err := podhelper.HasNoBadLogs(check.client, &pods.Items[i], "osm-controller"); err != nil {
			return err // TODO since we can have multiple osm-controller pods, should we return err on the first controller with bad logs?
		}
	}

	return nil
}

// Suggestion implements common.Runnable.
func (check NoBadOsmControllerLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadOsmControllerLogsCheck) FixIt() error {
	panic("implement me")
}
