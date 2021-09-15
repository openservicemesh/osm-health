package namespace

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm/pkg/constants"
)

const enabled = "enabled"

// Verify interface compliance
var _ runner.Runnable = (*SidecarInjectionCheck)(nil)

// SidecarInjectionCheck implements common.Runnable
type SidecarInjectionCheck struct {
	client    kubernetes.Interface
	namespace string
}

// Suggestion implements common.Runnable
func (check SidecarInjectionCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (check SidecarInjectionCheck) FixIt() error {
	panic("implement me")
}

// NewSidecarInjectionCheck creates a SidecarInjectionCheck which checks whether a namespace is enabled for sidecar injection.
func NewSidecarInjectionCheck(client kubernetes.Interface, namespace string) SidecarInjectionCheck {
	return SidecarInjectionCheck{
		client:    client,
		namespace: namespace,
	}
}

// Description implements common.Runnable
func (check SidecarInjectionCheck) Description() string {
	return fmt.Sprintf("Checking whether namespace %s is annotated for automatic sidecar injection", check.namespace)
}

// Run implements common.Runnable
func (check SidecarInjectionCheck) Run() outcomes.Outcome {
	annotations, err := getAnnotations(check.client, check.namespace)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return outcomes.Fail{Error: err}
	}

	annotationValue, ok := annotations[constants.SidecarInjectionAnnotation]
	isAnnotatedForInjection := ok && annotationValue == enabled

	if !isAnnotatedForInjection {
		return outcomes.Fail{Error: ErrNotAnnotatedForSidecarInjection}
	}

	return outcomes.Pass{}
}
