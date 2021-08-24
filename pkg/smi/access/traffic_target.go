package access

import (
	"context"
	"fmt"
	"strings"

	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha2"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha3"
	"github.com/openservicemesh/osm/pkg/configurator"
)

// Verify interface compliance
var _ common.Runnable = (*TrafficTargetCheck)(nil)

// TrafficTargetCheck implements common.Runnable
type TrafficTargetCheck struct {
	osmVersion   osm.ControllerVersion
	cfg          configurator.Configurator
	srcPod       *corev1.Pod
	dstPod       *corev1.Pod
	accessClient smiAccessClient.Interface
}

// NewTrafficTargetCheck creates a check of type TrafficTargetCheck which checks whether the src and dest pods are referenced as src and dest in a TrafficTarget (in that order)
func NewTrafficTargetCheck(osmVersion osm.ControllerVersion, osmConfigurator configurator.Configurator, srcPod *corev1.Pod, dstPod *corev1.Pod, smiAccessClient smiAccessClient.Interface) TrafficTargetCheck {
	return TrafficTargetCheck{
		osmVersion:   osmVersion,
		cfg:          osmConfigurator,
		srcPod:       srcPod,
		dstPod:       dstPod,
		accessClient: smiAccessClient,
	}
}

// Description implements common.Runnable
func (check TrafficTargetCheck) Description() string {
	return fmt.Sprintf(
		"Checking whether there is a Traffic Target with source pod %s and destination pod %s in namespace %s",
		check.srcPod.Name,
		check.dstPod.Name,
		check.dstPod.Namespace)
}

// Run implements common.Runnable
func (check TrafficTargetCheck) Run() outcomes.Outcome {
	// Check if permissive mode is enabled, in which case every meshed pod is allowed to communicate with each other
	if check.cfg.IsPermissiveTrafficPolicyMode() {
		return outcomes.DiagnosticOutcome{LongDiagnostics: "OSM is in permissive traffic policy modes -- all meshed pods can communicate and SMI access policies are not applicable"}
	}
	switch check.osmVersion {
	case "v0.5", "v0.6":
		return check.runForV1alpha2()
	case "v0.7", "v0.8", "v0.9":
		return check.runForV1alpha3()
	default:
		return outcomes.FailedOutcome{Error: fmt.Errorf(
			"OSM Controller version could not be mapped to a TrafficTarget version. Supported versions are v0.5 through v0.9")}
	}
}

func (check TrafficTargetCheck) runForV1alpha2() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha2().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.FailedOutcome{Error: err}
	}
	var matchingTargetNames []string
	for _, trafficTarget := range trafficTargets.Items {
		if v1alpha2.DoesTargetMatchPods(trafficTarget.Spec, check.srcPod, check.dstPod) {
			matchingTargetNames = append(matchingTargetNames, trafficTarget.Name)
		}
	}
	if len(matchingTargetNames) > 0 {
		return outcomes.DiagnosticOutcome{LongDiagnostics: fmt.Sprintf(
			"Pod '%s/%s' is allowed to communicate to pod '%s/%s' via SMI TrafficTarget policy/policies %s\n",
			check.srcPod.Namespace,
			check.srcPod.Name,
			check.dstPod.Namespace,
			check.dstPod.Name,
			strings.Join(matchingTargetNames, ", ")),
		}
	}
	return outcomes.DiagnosticOutcome{LongDiagnostics: fmt.Sprintf(
		"Pod '%s/%s' is not allowed to communicate to pod '%s/%s' via any SMI TrafficTarget policy\n",
		check.srcPod.Namespace,
		check.srcPod.Name,
		check.dstPod.Namespace,
		check.dstPod.Name)}
}

func (check TrafficTargetCheck) runForV1alpha3() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha3().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.FailedOutcome{Error: err}
	}
	var matchingTargetNames []string
	for _, trafficTarget := range trafficTargets.Items {
		if v1alpha3.DoesTargetMatchPods(trafficTarget.Spec, check.srcPod, check.dstPod) {
			matchingTargetNames = append(matchingTargetNames, trafficTarget.Name)
		}
	}
	if len(matchingTargetNames) > 0 {
		return outcomes.DiagnosticOutcome{LongDiagnostics: fmt.Sprintf(
			"Pod '%s/%s' is allowed to communicate to pod '%s/%s' via SMI TrafficTarget policy/policies %s\n",
			check.srcPod.Namespace,
			check.srcPod.Name,
			check.dstPod.Namespace,
			check.dstPod.Name,
			strings.Join(matchingTargetNames, ", "))}
	}
	return outcomes.DiagnosticOutcome{LongDiagnostics: fmt.Sprintf(
		"Pod '%s/%s' is not allowed to communicate to pod '%s/%s' via any SMI TrafficTarget policy\n",
		check.srcPod.Namespace,
		check.srcPod.Name,
		check.dstPod.Namespace,
		check.dstPod.Name)}
}

// Suggestion implements common.Runnable
func (check TrafficTargetCheck) Suggestion() string {
	return fmt.Sprintf(
		"Check that source and desintation pod are referred to in a TrafficTarget. To get relevant TrafficTargets, use: \"kubectl get traffictarget -n %s -o yaml\"",
		check.dstPod.Namespace)
}

// FixIt implements common.Runnable
func (check TrafficTargetCheck) FixIt() error {
	panic("implement me")
}
