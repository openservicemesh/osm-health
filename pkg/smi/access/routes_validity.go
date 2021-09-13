package access

import (
	"context"
	"fmt"

	"github.com/openservicemesh/osm-health/pkg/runner"

	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/smi"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha2"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha3"
	"github.com/openservicemesh/osm/pkg/configurator"
)

// Verify interface compliance
var _ runner.Runnable = (*RoutesValidityCheck)(nil)

// RoutesValidityCheck implements common.Runnable
type RoutesValidityCheck struct {
	osmVersion   osm.ControllerVersion
	cfg          configurator.Configurator
	srcPod       *corev1.Pod
	dstPod       *corev1.Pod
	accessClient smiAccessClient.Interface
}

// NewRoutesValidityCheck returns a check of type RoutesValidityCheck which checks whether TrafficTargets matching the src and dest pods have supported routes
func NewRoutesValidityCheck(osmVersion osm.ControllerVersion, osmConfigurator configurator.Configurator, srcPod *corev1.Pod, dstPod *corev1.Pod, smiAccessClient smiAccessClient.Interface) RoutesValidityCheck {
	return RoutesValidityCheck{
		osmVersion:   osmVersion,
		cfg:          osmConfigurator,
		srcPod:       srcPod,
		dstPod:       dstPod,
		accessClient: smiAccessClient,
	}
}

// Description implements common.Runnable
func (check RoutesValidityCheck) Description() string {
	return fmt.Sprintf(
		"Checking whether Traffic Targets in namespace %s with source pod %s and destination pod %s have valid routes (Kind: %s or %s)",
		check.dstPod.Namespace,
		check.srcPod.Name,
		check.dstPod.Name,
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind)
}

// Run implements common.Runnable
func (check RoutesValidityCheck) Run() outcomes.Outcome {
	// Check if permissive mode is enabled, in which case every meshed pod is allowed to communicate with each other
	if check.cfg.IsPermissiveTrafficPolicyMode() {
		return outcomes.Info{Diagnostics: "OSM is in permissive traffic policy modes -- all meshed pods can communicate and SMI access policies are not applicable"}
	}
	switch osm.SupportedTrafficTarget[check.osmVersion] {
	case osm.V1Alpha2:
		return check.runForTrafficTargetV1alpha2()
	case osm.V1Alpha3:
		return check.runForTrafficTargetV1alpha3()
	default:
		return outcomes.Fail{Error: fmt.Errorf(
			"OSM Controller version could not be mapped to a TrafficTarget version. Supported versions are v0.5 through v0.9")}
	}
}

func (check RoutesValidityCheck) runForTrafficTargetV1alpha2() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha2().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.Fail{Error: err}
	}
	unsupportedRouteTargets := map[string]string{}
	var foundMatchingTarget bool
	for _, trafficTarget := range trafficTargets.Items {
		spec := trafficTarget.Spec
		if !v1alpha2.DoesTargetMatchPods(spec, check.srcPod, check.dstPod) {
			continue
		}
		foundMatchingTarget = true
		for _, rule := range spec.Rules {
			kind := rule.Kind
			err = isTrafficTargetRouteKindSupported(kind, check.osmVersion)
			if err != nil {
				unsupportedRouteTargets[trafficTarget.Name] = kind
			}
		}
	}
	if !foundMatchingTarget {
		return outcomes.Info{Diagnostics: fmt.Sprintf(
			"No applicable Traffic Targets in namespace %s to check routes for",
			check.dstPod.Namespace)}
	}
	if len(unsupportedRouteTargets) > 0 {
		errorString := check.newErrorMessage(unsupportedRouteTargets)
		return outcomes.Fail{Error: fmt.Errorf(errorString)}
	}
	return outcomes.Pass{}
}

func (check RoutesValidityCheck) runForTrafficTargetV1alpha3() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha3().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.Fail{Error: err}
	}
	unsupportedRouteTargets := map[string]string{}
	var foundMatchingTarget bool
	for _, trafficTarget := range trafficTargets.Items {
		spec := trafficTarget.Spec
		if !v1alpha3.DoesTargetMatchPods(spec, check.srcPod, check.dstPod) {
			continue
		}
		foundMatchingTarget = true
		for _, rule := range spec.Rules {
			kind := rule.Kind
			err = isTrafficTargetRouteKindSupported(kind, check.osmVersion)
			if err != nil {
				unsupportedRouteTargets[trafficTarget.Name] = kind
			}
		}
	}
	if !foundMatchingTarget {
		return outcomes.Info{Diagnostics: fmt.Sprintf(
			"No applicable Traffic Targets in namespace %s to check routes for",
			check.dstPod.Namespace)}
	}
	if len(unsupportedRouteTargets) > 0 {
		errorString := check.newErrorMessage(unsupportedRouteTargets)
		return outcomes.Fail{Error: fmt.Errorf(errorString)}
	}
	return outcomes.Pass{}
}

func (*RoutesValidityCheck) newErrorMessage(targetToKindMap map[string]string) string {
	errorString := fmt.Sprintf(
		"Expected routes of kind %s or %s, found the following TrafficTargets with unsupported routes: \n",
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind)
	for target, kind := range targetToKindMap {
		errorString = fmt.Sprintf("%s %s: %s\n", errorString, target, kind)
	}
	return errorString
}

// Suggestion implements common.Runnable
func (check RoutesValidityCheck) Suggestion() string {
	return fmt.Sprintf(
		"Check that TrafficTargets routes are of kind %s or %s. To get relevant TrafficTargets, use: \"kubectl get traffictarget -n %s -o yaml\"",
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
		check.dstPod.Namespace)
}

// FixIt implements common.Runnable
func (check RoutesValidityCheck) FixIt() error {
	panic("implement me")
}
