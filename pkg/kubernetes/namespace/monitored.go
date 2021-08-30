package namespace

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
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

// NewMonitoredCheck creates a MonitoredCheck which checks whether a namespace is monitored by certain OSM Controller.
func NewMonitoredCheck(client kubernetes.Interface, namespace string, meshName common.MeshName) MonitoredCheck {
	return MonitoredCheck{
		client:    client,
		namespace: namespace,
		meshName:  meshName,
	}
}

// Description implements common.Runnable
func (check MonitoredCheck) Description() string {
	return fmt.Sprintf("Checking whether namespace %s is monitored by OSM %s", check.namespace, check.meshName)
}

// Run implements common.Runnable
func (check MonitoredCheck) Run() outcomes.Outcome {
	labels, err := getLabels(check.client, check.namespace)
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	labelValue, ok := labels[constants.OSMKubeResourceMonitorAnnotation]
	isMonitoredByController := ok && labelValue == check.meshName.String()

	if !isMonitoredByController {
		return outcomes.Fail{Error: ErrNotMonitoredByOSMController}
	}

	return outcomes.Pass{}
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

// NewNamespacesInSameMeshCheck creates a SidecarInjectionCheck which checks whether two pods are in the same mesh
func NewNamespacesInSameMeshCheck(client kubernetes.Interface, namespaceA string, namespaceB string) NamespacesInSameMeshCheck {
	return NamespacesInSameMeshCheck{
		client:     client,
		namespaceA: namespaceA,
		namespaceB: namespaceB,
	}
}

// Description implements common.Runnable
func (check NamespacesInSameMeshCheck) Description() string {
	return fmt.Sprintf("Checking whether namespace %s and namespace %s are monitored by the same mesh", check.namespaceA, check.namespaceB)
}

// Run implements common.Runnable
func (check NamespacesInSameMeshCheck) Run() outcomes.Outcome {
	labelsA, err := getLabels(check.client, check.namespaceA)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	meshNameA, labelExistsA := labelsA[constants.OSMKubeResourceMonitorAnnotation]

	labelsB, err := getLabels(check.client, check.namespaceB)
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	meshNameB, labelExistsB := labelsB[constants.OSMKubeResourceMonitorAnnotation]
	if !labelExistsA || !labelExistsB {
		return outcomes.Fail{Error: ErrNotMonitoredByOSMController}
	}
	if meshNameA != meshNameB {
		return outcomes.Fail{Error: ErrNamespacesNotInSameMesh}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (check NamespacesInSameMeshCheck) Suggestion() string {
	return fmt.Sprintf("Verify which mesh each namespace is monitored by. Try: \"kubectl get namespace -n %s -o json | jq '.items[0].metadata.labels'\"", check.namespaceA)
}

// FixIt implements common.Runnable
func (check NamespacesInSameMeshCheck) FixIt() error {
	panic("implement me")
}
