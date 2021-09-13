package osm

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/utils"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/k8s"
)

// Verify interface compliance
var _ runner.Runnable = (*ControllerHTTPServerEndpointsCheck)(nil)

// ControllerHTTPServerEndpointsCheck implements common.Runnable
type ControllerHTTPServerEndpointsCheck struct {
	client                   kubernetes.Interface
	osmControlPlaneNamespace common.MeshNamespace
	localPort                uint16
	actionConfig             *action.Configuration
}

// HasValidInfoFromControllerHTTPServerEndpointsCheck checks whether the osm-controller's http server endpoints return valid information.
func HasValidInfoFromControllerHTTPServerEndpointsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace, localPort uint16, actionConfig *action.Configuration) ControllerHTTPServerEndpointsCheck {
	return ControllerHTTPServerEndpointsCheck{
		client:                   client,
		osmControlPlaneNamespace: osmControlPlaneNamespace,
		localPort:                localPort,
		actionConfig:             actionConfig,
	}
}

// Description implements common.Runnable
func (check ControllerHTTPServerEndpointsCheck) Description() string {
	return "Checking whether osm-controller http server endpoints return valid information"
}

// GetOSMControllerPods TODO remove once osm PR for exporting GetOSMControllerPods through osm/pkg/cli is merged.
func GetOSMControllerPods(clientSet kubernetes.Interface, ns string) *corev1.PodList {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMControllerName}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	podList, _ := clientSet.CoreV1().Pods(ns).List(context.TODO(), listOptions)
	return podList
}

// Run implements common.Runnable
func (check ControllerHTTPServerEndpointsCheck) Run() outcomes.Outcome {
	// TODO replace with cli.GetOSMControllerPods(...) once osm repo PR is merged
	controllerPods := GetOSMControllerPods(check.client, check.osmControlPlaneNamespace.String())

	anyControllerPodsExist := false
	for _, controllerPod := range controllerPods.Items {
		anyControllerPodsExist = true

		conf, err := check.actionConfig.RESTClientGetter.ToRESTConfig()
		if err != nil {
			return outcomes.Fail{Error: errors.Errorf("failed to get REST config from Helm %s", err)}
		}
		dialer, err := k8s.DialerToPod(conf, check.client, controllerPod.Name, controllerPod.Namespace)
		if err != nil {
			return outcomes.Fail{Error: errors.Errorf("error setting up port forwarding: %s", err)}
		}
		portForwarder, err := k8s.NewPortForwarder(dialer, fmt.Sprintf("%d:%d", check.localPort, constants.OSMHTTPServerPort))
		if err != nil {
			return outcomes.Fail{Error: errors.Errorf("error setting up port forwarding: %s", err)}
		}

		err = portForwarder.Start(func(pf *k8s.PortForwarder) error {
			defer pf.Stop()
			controllerHTTPServerURL := fmt.Sprintf("http://localhost:%d", check.localPort)

			err = checkControllerHealthReadiness(controllerHTTPServerURL)
			if err != nil {
				return err
			}

			err = checkControllerHealthLiveness(controllerHTTPServerURL)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return outcomes.Fail{Error: err}
		}
	}

	if !anyControllerPodsExist {
		return outcomes.Fail{Error: ErrorNoControllerPodsExistInNamespace}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable.
func (check ControllerHTTPServerEndpointsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check ControllerHTTPServerEndpointsCheck) FixIt() error {
	panic("implement me")
}

func checkControllerHealthReadiness(controllerHTTPServerURL string) error {
	// TODO replace "/health/ready" with constants.HTTPServer... from osm repo once PR over there is merged.
	url := fmt.Sprintf("%s%s", controllerHTTPServerURL, "/health/ready")
	respBody, err := utils.GetResponseBody(url)
	if err != nil {
		return errors.Errorf("osm-controller health readiness check failed: %s", err)
	}

	if respBody != "Service is ready" {
		return ErrorControllerNotReady
	}
	return nil
}

func checkControllerHealthLiveness(controllerHTTPServerURL string) error {
	// TODO replace "/health/alive" with constants.HTTPServer... from osm repo once PR over there is merged.
	url := fmt.Sprintf("%s%s", controllerHTTPServerURL, "/health/alive")
	respBody, err := utils.GetResponseBody(url)
	if err != nil {
		return errors.Errorf("osm-controller health liveness check failed: %s", err)
	}

	if respBody != "Service is alive" {
		return ErrorControllerNotAlive
	}
	return nil
}
