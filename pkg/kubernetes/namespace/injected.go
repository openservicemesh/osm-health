package namespace

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

const enabled = "enabled"

// SidecarInjectionCheck implements common.Runnable
type SidecarInjectionCheck struct {
	client    kubernetes.Interface
	namespace string
}

// IsInjectEnabled checks whether a namespace is enabled for sidecar injection.
func IsInjectEnabled(client kubernetes.Interface, namespace string) common.Runnable {
	return isInjectEnabled(client, namespace)
}

func isInjectEnabled(client kubernetes.Interface, namespace string) SidecarInjectionCheck {
	return SidecarInjectionCheck{
		client:    client,
		namespace: namespace,
	}
}

// Info implements common.Runnable
func (check SidecarInjectionCheck) Info() string {
	return fmt.Sprintf("Checking whether namespace %s is annotated for automatic sidecar injection", check.namespace)
}

// Run implements common.Runnable
func (check SidecarInjectionCheck) Run() error {
	annotations, err := getAnnotations(check.client, check.namespace)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	annotationValue, ok := annotations[constants.SidecarInjectionAnnotation]
	isAnnotatedForInjection := ok && annotationValue == enabled

	if !isAnnotatedForInjection {
		return ErrNotAnnotatedForSidecarInjection
	}

	return nil
}
