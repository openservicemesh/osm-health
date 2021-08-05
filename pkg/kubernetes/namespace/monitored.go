package namespace

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

// Verify interface compliance
var _ common.Runnable = (*MonitoredCheck)(nil)

// MonitoredCheck implements common.Runnable
type MonitoredCheck struct {
	client    kubernetes.Interface
	namespace string
	meshName  common.MeshName
}

// IsMonitoredBy checks whether a namespace is monitored by certain OSM Controller.
func IsMonitoredBy(client kubernetes.Interface, namespace string, meshName common.MeshName) MonitoredCheck {
	return MonitoredCheck{
		client:    client,
		namespace: namespace,
		meshName:  meshName,
	}
}

// Info implements common.Runnable
func (check MonitoredCheck) Info() string {
	return fmt.Sprintf("Checking whether namespace %s is monitored by OSM %s", check.namespace, check.meshName)
}

// Run implements common.Runnable
func (check MonitoredCheck) Run() error {
	labels, err := getLabels(check.client, check.namespace)
	if err != nil {
		return err
	}

	labelValue, ok := labels[constants.OSMKubeResourceMonitorAnnotation]
	isMonitoredByController := ok && labelValue == check.meshName.String()

	if !isMonitoredByController {
		return ErrNotMonitoredByOSMController
	}

	return nil
}

// Suggestion implements common.Runnable
func (check MonitoredCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (check MonitoredCheck) FixIt() error {
	panic("implement me")
}

// Verify interface compliance
var _ common.Runnable = (*NamespacesInSameMeshCheck)(nil)

// NamespacesInSameMeshCheck implements common.Runnable
type NamespacesInSameMeshCheck struct {
	client     kubernetes.Interface
	namespaceA string
	namespaceB string
}

// AreNamespacesInSameMesh checks whether two pods are in the same mesh
func AreNamespacesInSameMesh(client kubernetes.Interface, namespaceA string, namespaceB string) NamespacesInSameMeshCheck {
	return NamespacesInSameMeshCheck{
		client:     client,
		namespaceA: namespaceA,
		namespaceB: namespaceB,
	}
}

// Info implements common.Runnable
func (check NamespacesInSameMeshCheck) Info() string {
	return fmt.Sprintf("Checking whether namespace %s and namespace %s are monitored by the same mesh", check.namespaceA, check.namespaceB)
}

// Run implements common.Runnable
func (check NamespacesInSameMeshCheck) Run() error {
	labelsA, err := getLabels(check.client, check.namespaceA)
	if err != nil {
		return err
	}
	meshNameA, labelExistsA := labelsA[constants.OSMKubeResourceMonitorAnnotation]

	labelsB, err := getLabels(check.client, check.namespaceB)
	if err != nil {
		return err
	}

	meshNameB, labelExistsB := labelsB[constants.OSMKubeResourceMonitorAnnotation]
	if !labelExistsA || !labelExistsB {
		return ErrNotMonitoredByOSMController
	}
	if meshNameA != meshNameB {
		return ErrNamespacesNotInSameMesh
	}
	return nil
}

// Suggestion implements common.Runnable
func (check NamespacesInSameMeshCheck) Suggestion() string {
	return fmt.Sprintf("Verify which mesh each namespace is monitored by. Try: \"kubectl get namespace -n %s -o json | jq '.items[0].metadata.labels'\"", check.namespaceA)
}

// FixIt implements common.Runnable
func (check NamespacesInSameMeshCheck) FixIt() error {
	panic("implement me")
}
