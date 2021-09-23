package access

import (
	"context"
	"fmt"
	"strings"

	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
	smiSpecClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/specs/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/osm/version"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha2"
	"github.com/openservicemesh/osm-health/pkg/smi/access/v1alpha3"
	"github.com/openservicemesh/osm/pkg/configurator"
)

// Verify interface compliance
var _ runner.Runnable = (*RoutesExistenceCheck)(nil)

// RoutesExistenceCheck implements common.Runnable
type RoutesExistenceCheck struct {
	osmVersion   version.ControllerVersion
	cfg          configurator.Configurator
	srcPod       *corev1.Pod
	dstPod       *corev1.Pod
	accessClient smiAccessClient.Interface
	specClient   smiSpecClient.Interface
}

// NewRoutesExistenceCheck checks whether routes referenced by TrafficTargets matching the src and dest pods exist in the cluster
func NewRoutesExistenceCheck(osmVersion version.ControllerVersion, osmConfigurator configurator.Configurator, srcPod *corev1.Pod, dstPod *corev1.Pod, smiAccessClient smiAccessClient.Interface, smiSpecClient smiSpecClient.Interface) RoutesExistenceCheck {
	return RoutesExistenceCheck{
		osmVersion:   osmVersion,
		cfg:          osmConfigurator,
		srcPod:       srcPod,
		dstPod:       dstPod,
		accessClient: smiAccessClient,
		specClient:   smiSpecClient,
	}
}

// Description implements common.Runnable
func (check RoutesExistenceCheck) Description() string {
	return fmt.Sprintf("Checking whether routes referenced by above matched TrafficTargets exist in namespace %s", check.dstPod.Namespace)
}

// Run implements common.Runnable
func (check RoutesExistenceCheck) Run() outcomes.Outcome {
	// Check if permissive mode is enabled, in which case every meshed pod is allowed to communicate with each other
	if check.cfg.IsPermissiveTrafficPolicyMode() {
		return outcomes.Info{Diagnostics: "OSM is in permissive traffic policy modes -- all meshed pods can communicate and SMI access policies are not applicable"}
	}
	switch version.SupportedTrafficTarget[check.osmVersion] {
	case version.V1Alpha2:
		return check.runForTrafficTargetV1alpha2()
	case version.V1Alpha3:
		return check.runForTrafficTargetV1alpha3()
	default:
		return outcomes.Fail{Error: fmt.Errorf(
			"OSM Controller version could not be mapped to a TrafficTarget version. Supported versions are v0.5 through v0.9")}
	}
}

func (check RoutesExistenceCheck) runForTrafficTargetV1alpha2() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha2().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.Fail{Error: err}
	}
	existingRoutes, err := v1alpha2.GetExistingRouteNames(check.specClient, check.dstPod.Namespace)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if existingRoutes.Cardinality() == 0 {
		return outcomes.Fail{Error: fmt.Errorf("No HTTPRouteGroups or TCPRoutes exist in namespace %s", check.dstPod.Namespace)}
	}

	var missingRoutes []string
	var foundMatchingTarget bool
	for _, trafficTarget := range trafficTargets.Items {
		spec := trafficTarget.Spec
		if !v1alpha2.DoesTargetMatchPods(spec, check.srcPod, check.dstPod) {
			continue
		}
		foundMatchingTarget = true
		for _, rule := range spec.Rules {
			err = isTrafficTargetRouteKindSupported(rule.Kind, check.osmVersion)
			if err == nil && !(existingRoutes.Contains(rule.Name)) {
				missingRoutes = append(missingRoutes, rule.Name)
			}
		}
	}
	if !foundMatchingTarget {
		return outcomes.Info{Diagnostics: fmt.Sprintf("No applicable Traffic Targets in namespace %s to check routes for", check.dstPod.Namespace)}
	}
	if len(missingRoutes) > 0 {
		return outcomes.Fail{Error: fmt.Errorf("The following routes could not be found in the cluster: %s", strings.Join(missingRoutes, ", "))}
	}
	return outcomes.Pass{}
}

func (check RoutesExistenceCheck) runForTrafficTargetV1alpha3() outcomes.Outcome {
	trafficTargets, err := check.accessClient.AccessV1alpha3().TrafficTargets(check.dstPod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting TrafficTargets for namespace %s", check.dstPod.Namespace)
		return outcomes.Fail{Error: err}
	}
	existingRoutes, err := v1alpha3.GetExistingRouteNames(check.specClient, check.dstPod.Namespace)
	if err != nil {
		return outcomes.Fail{Error: err}
	}
	if existingRoutes.Cardinality() == 0 {
		return outcomes.Fail{Error: fmt.Errorf("No HTTPRouteGroups or TCPRoutes exist in namespace %s", check.dstPod.Namespace)}
	}

	var missingRoutes []string
	var foundMatchingTarget bool
	for _, trafficTarget := range trafficTargets.Items {
		spec := trafficTarget.Spec
		if !v1alpha3.DoesTargetMatchPods(spec, check.srcPod, check.dstPod) {
			continue
		}
		foundMatchingTarget = true
		for _, rule := range spec.Rules {
			err = isTrafficTargetRouteKindSupported(rule.Kind, check.osmVersion)
			if err == nil && !(existingRoutes.Contains(rule.Name)) {
				missingRoutes = append(missingRoutes, rule.Name)
			}
		}
	}
	if !foundMatchingTarget {
		return outcomes.Info{Diagnostics: fmt.Sprintf("No applicable Traffic Targets in namespace %s to check routes for", check.dstPod.Namespace)}
	}
	if len(missingRoutes) > 0 {
		return outcomes.Fail{Error: fmt.Errorf("The following routes could not be found in the cluster: %s", strings.Join(missingRoutes, ", "))}
	}
	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (check RoutesExistenceCheck) Suggestion() string {
	return fmt.Sprintf("Check that routes referenced by the TrafficTargets exist in the cluster. Use: \"kubectl get httproutegroups -n %s\" and \"kubectl get tcproutes -n %s\"", check.dstPod.Namespace, check.dstPod.Namespace)
}

// FixIt implements common.Runnable
func (check RoutesExistenceCheck) FixIt() error {
	panic("implement me")
}
