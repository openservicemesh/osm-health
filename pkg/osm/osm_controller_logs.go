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

// Description implements common.Runnable
func (check NoBadOsmControllerLogsCheck) Description() string {
	return fmt.Sprintf("Checking whether namespace %s has bad (fatal/error/warning/fail) logs in osm controller pods", check.osmControlPlaneNamespace)
}

// Run implements common.Runnable
func (check NoBadOsmControllerLogsCheck) Run() outcomes.Outcome {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMControllerName}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	pods, err := check.client.CoreV1().Pods(check.osmControlPlaneNamespace).List(context.TODO(), listOptions)
	if err != nil {
		return outcomes.FailedOutcome{Error: fmt.Errorf("unable to list %s pods in namespace %s", constants.OSMControllerName, check.osmControlPlaneNamespace)}
	}

	osmControllerErrCount := 0
	for i := range pods.Items {
		if err := podhelper.HasNoBadLogs(check.client, &pods.Items[i], "osm-controller").GetError(); err != nil {
			osmControllerErrCount++
			log.Error().Err(err)
		}
	}

	if osmControllerErrCount != 0 {
		return outcomes.FailedOutcome{Error: errors.Errorf("%s pods in namespace %s have %d errors", constants.OSMControllerName, check.osmControlPlaneNamespace, osmControllerErrCount)}
	}

	return outcomes.SuccessfulOutcomeWithoutDiagnostics{}
}

// Suggestion implements common.Runnable.
func (check NoBadOsmControllerLogsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check NoBadOsmControllerLogsCheck) FixIt() error {
	panic("implement me")
}
