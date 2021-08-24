package v1alpha3

import (
	"context"

	mapset "github.com/deckarep/golang-set"
	accessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha3"
	smiSpecClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/specs/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openservicemesh/osm/pkg/cli"
)

// DoesTargetMatchPods checks whether a given TrafficTarget has dstPod as its destination as dstPod and srcPod as an allowed source to this destination
func DoesTargetMatchPods(spec accessClient.TrafficTargetSpec, srcPod *corev1.Pod, dstPod *corev1.Pod) bool {
	return cli.DoesTargetRefDstPod(spec, dstPod) && cli.DoesTargetRefSrcPod(spec, srcPod)
}

// GetExistingRouteNames returns the names of HTTPRouteGroups and TCPRoutes that exist in the cluster
func GetExistingRouteNames(specClient smiSpecClient.Interface, namespace string) (mapset.Set, error) {
	routes := mapset.NewSet()
	httpRouteGroups, err := specClient.SpecsV1alpha4().HTTPRouteGroups(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting HTTPRouteGroups for namespace %s", namespace)
	}
	for _, httpRouteGroup := range httpRouteGroups.Items {
		routes.Add(httpRouteGroup.Name)
	}
	tcpRoutes, err := specClient.SpecsV1alpha4().TCPRoutes(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting HTTPRouteGroups for namespace %s", namespace)
	}
	for _, tcpRoute := range tcpRoutes.Items {
		routes.Add(tcpRoute.Name)
	}
	return routes, err
}
