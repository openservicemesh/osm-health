package namespace

import (
	"fmt"

	kubernetes2 "github.com/openservicemesh/osm-health/pkg/kubernetes"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm/pkg/constants"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// MonitoredCheck implements common.Runnable
type MonitoredCheck struct {
	client    kubernetes.Interface
	namespace kubernetes2.Namespace
	meshName  common.MeshName
}

// IsMonitoredBy checks whether a namespace is monitored by certain OSM Controller.
func IsMonitoredBy(client kubernetes.Interface, namespace kubernetes2.Namespace, meshName common.MeshName) common.Runnable {
	return isMonitoredBy(client, namespace, meshName)
}

func isMonitoredBy(client kubernetes.Interface, namespace kubernetes2.Namespace, meshName common.MeshName) MonitoredCheck {
	return MonitoredCheck{
		client:    client,
		namespace: namespace,
		meshName:  meshName,
	}
}

// Info implements common.Runnable
func (mc MonitoredCheck) Info() string {
	return fmt.Sprintf("Checking whether namespace %s is monitored by OSM %s", mc.namespace, mc.meshName)
}

// Run implements common.Runnable
func (mc MonitoredCheck) Run() error {
	labels, err := getLabels(mc.client, mc.namespace)
	if err != nil {
		return err
	}

	labelValue, ok := labels[constants.OSMKubeResourceMonitorAnnotation]
	isMonitoredByController := ok && labelValue == mc.meshName.String()

	if !isMonitoredByController {
		return ErrNotMonitoredByOSMController
	}

	return nil
}
