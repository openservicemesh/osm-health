package controller

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/utils"
	"github.com/openservicemesh/osm/pkg/constants"
	httpserverconstants "github.com/openservicemesh/osm/pkg/httpserver/constants"
	"github.com/openservicemesh/osm/pkg/k8s"
)

// Verify interface compliance
var _ runner.Runnable = (*HTTPServerEndpointsCheck)(nil)

// HTTPServerEndpointsCheck implements common.Runnable
type HTTPServerEndpointsCheck struct {
	client                   kubernetes.Interface
	osmControlPlaneNamespace common.MeshNamespace
	localPort                uint16
	actionConfig             *action.Configuration
}

// HasValidInfoFromControllerHTTPServerEndpointsCheck checks whether the osm-controller's http server endpoints return valid information.
func HasValidInfoFromControllerHTTPServerEndpointsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace, localPort uint16, actionConfig *action.Configuration) HTTPServerEndpointsCheck {
	return HTTPServerEndpointsCheck{
		client:                   client,
		osmControlPlaneNamespace: osmControlPlaneNamespace,
		localPort:                localPort,
		actionConfig:             actionConfig,
	}
}

// Description implements common.Runnable
func (check HTTPServerEndpointsCheck) Description() string {
	return "Checking whether osm-controller http server endpoints return valid information"
}

// Run implements common.Runnable
func (check HTTPServerEndpointsCheck) Run() outcomes.Outcome {
	controllerPods := k8s.GetOSMControllerPods(check.client, check.osmControlPlaneNamespace.String())

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
func (check HTTPServerEndpointsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check HTTPServerEndpointsCheck) FixIt() error {
	panic("implement me")
}

func checkControllerHealthReadiness(controllerHTTPServerURL string) error {
	url := fmt.Sprintf("%s%s", controllerHTTPServerURL, httpserverconstants.HealthReadinessPath)
	if err := utils.IsHTTPResponseCodeEquals(url, http.StatusOK); err != nil {
		return errors.Errorf("osm-controller health readiness check failed: %s", err)
	}
	return nil
}

func checkControllerHealthLiveness(controllerHTTPServerURL string) error {
	url := fmt.Sprintf("%s%s", controllerHTTPServerURL, httpserverconstants.HealthLivenessPath)
	if err := utils.IsHTTPResponseCodeEquals(url, http.StatusOK); err != nil {
		return errors.Errorf("osm-controller health liveness check failed: %s", err)
	}
	return nil
}
