package v1alpha3

import (
	accessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha3"
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm/pkg/cli"
)

// DoesTargetMatchPods checks whether a given TrafficTarget has dstPod as its destination as dstPod and srcPod as an allowed source to this destination
func DoesTargetMatchPods(spec accessClient.TrafficTargetSpec, srcPod *corev1.Pod, dstPod *corev1.Pod) bool {
	return cli.DoesTargetRefDstPod(spec, dstPod) && cli.DoesTargetRefSrcPod(spec, srcPod)
}
