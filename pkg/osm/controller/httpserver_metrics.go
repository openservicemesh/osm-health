package controller

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	osmutils "github.com/openservicemesh/osm-health/pkg/osm/utils"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/utils"
	"github.com/openservicemesh/osm/pkg/constants"
	httpserverconstants "github.com/openservicemesh/osm/pkg/httpserver/constants"
	"github.com/openservicemesh/osm/pkg/k8s"
)

// Verify interface compliance
var _ runner.Runnable = (*HTTPServerProxyConnectionMetricsCheck)(nil)

// HTTPServerProxyConnectionMetricsCheck implements common.Runnable
type HTTPServerProxyConnectionMetricsCheck struct {
	client                   kubernetes.Interface
	osmControlPlaneNamespace common.MeshNamespace
	controllerPods           *corev1.PodList
	localPort                uint16
	actionConfig             *action.Configuration
}

// NewHTTPServerProxyConnectionMetricsCheck checks whether the osm-controller's http server returns valid metrics for proxy connection count.
func NewHTTPServerProxyConnectionMetricsCheck(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace, controllerPods *corev1.PodList, localPort uint16, actionConfig *action.Configuration) HTTPServerProxyConnectionMetricsCheck {
	return HTTPServerProxyConnectionMetricsCheck{
		client:                   client,
		osmControlPlaneNamespace: osmControlPlaneNamespace,
		controllerPods:           controllerPods,
		localPort:                localPort,
		actionConfig:             actionConfig,
	}
}

// Description implements common.Runnable
func (check HTTPServerProxyConnectionMetricsCheck) Description() string {
	return "Checking whether the osm-controller's http server returns valid metrics for proxy connection count"
}

// Run implements common.Runnable
func (check HTTPServerProxyConnectionMetricsCheck) Run() outcomes.Outcome {
	anyControllerPodsExist := false
	for _, controllerPod := range check.controllerPods.Items {
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

			err = checkControllerProxyConnectionMetrics(check.client, controllerHTTPServerURL, check.osmControlPlaneNamespace)
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
func (check HTTPServerProxyConnectionMetricsCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable.
func (check HTTPServerProxyConnectionMetricsCheck) FixIt() error {
	panic("implement me")
}

func checkControllerProxyConnectionMetrics(client kubernetes.Interface, controllerHTTPServerURL string, osmControlPlaneNamespace common.MeshNamespace) error {
	url := fmt.Sprintf("%s%s", controllerHTTPServerURL, httpserverconstants.MetricsPath)
	metricsRespBody, err := utils.GetResponseBody(url)
	if err != nil {
		return errors.Errorf("osm-controller metrics check failed: %s", err)
	}

	monitoredNamespaces, err := osmutils.GetMonitoredNamespaces(client, osmControlPlaneNamespace)
	if err != nil {
		return errors.Errorf("osm-controller metrics check failed: %s", err)
	}

	// TODO - clarify if it is possible for a pod in a monitored namespace to NOT be a part of the mesh (have no proxy OR does not contribute to osm_proxy_connect_count)
	// TODO - clarify if it is still needed to check pod annotations and labels
	// TODO - it seems like when a namespace is ignored (through `osm namespace ignore ...`), the osm_proxy_connect_count is NOT decreased.
	// TODO - should we check for the metrics enabled annotation/label?
	totalMeshMonitoredPodsCount := 0
	for _, ns := range monitoredNamespaces.Items {
		pods, err := client.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return errors.Errorf("unable to list pods in monitored namespace %s", ns.Name)
		}
		totalMeshMonitoredPodsCount += len(pods.Items)
	}

	err = checkProxyConnectCount(totalMeshMonitoredPodsCount, metricsRespBody)
	if err != nil {
		return errors.Errorf("osm-controller metrics check failed: %s", err)
	}

	return nil
}

// checkProxyConnectCount checks whether the HTTP Server /metrics response body contains
// the correct value for the osm_proxy_connect_count metric.
func checkProxyConnectCount(expectedProxyConnectCount int, httpServerMetricsRespBody string) error {
	var osmProxyConnectCountMetricID = "osm_proxy_connect_count"
	// first capture group of the regex "(\\d+)" represents the number of proxies connected to OSM controller
	// Sample osm_proxy_connect_count metric returned from /metrics:
	//		# HELP osm_proxy_connect_count Represents the number of proxies connected to OSM controller
	//		# TYPE osm_proxy_connect_count gauge
	//		osm_proxy_connect_count 6
	r := regexp.MustCompile(fmt.Sprintf("%s\\s+(\\d+)", osmProxyConnectCountMetricID))
	match := r.FindStringSubmatch(httpServerMetricsRespBody)
	if len(match) != 2 {
		return errors.Errorf("missing or invalid %s metric in HTTP server metrics response", osmProxyConnectCountMetricID)
	}

	actualProxyConnectCount, err := strconv.Atoi(match[1])
	if err != nil {
		return errors.Errorf("invalid %s metric in HTTP server metrics response: %s", osmProxyConnectCountMetricID, err.Error())
	}

	if expectedProxyConnectCount != actualProxyConnectCount {
		return errors.Errorf("incorrect %s metric: expected %d but http server metrics returned %d",
			osmProxyConnectCountMetricID, expectedProxyConnectCount, actualProxyConnectCount)
	}

	return nil
}
