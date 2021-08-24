package v1alpha2

import (
	accessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha2"
	corev1 "k8s.io/api/core/v1"
)

const (
	serviceAccountKind = "ServiceAccount"
)

// DoesTargetMatchPods checks whether a given TrafficTarget has dstPod as its destination as dstPod and srcPod as an allowed source to this destination
func DoesTargetMatchPods(spec accessClient.TrafficTargetSpec, srcPod *corev1.Pod, dstPod *corev1.Pod) bool {
	return doesTargetRefDstPod(spec, dstPod) && doesTargetRefSrcPod(spec, srcPod)
}

// doesTargetRefDstPod checks whether the TrafficTarget spec refers to the destination pod's service account
func doesTargetRefDstPod(spec accessClient.TrafficTargetSpec, dstPod *corev1.Pod) bool {
	if spec.Destination.Kind != serviceAccountKind {
		return false
	}
	return spec.Destination.Name == dstPod.Spec.ServiceAccountName && spec.Destination.Namespace == dstPod.Namespace
}

// doesTargetRefSrcPod checks whether the TrafficTarget spec refers to the source pod's service account
func doesTargetRefSrcPod(spec accessClient.TrafficTargetSpec, srcPod *corev1.Pod) bool {
	for _, source := range spec.Sources {
		if source.Kind != serviceAccountKind {
			continue
		}
		if source.Name == srcPod.Spec.ServiceAccountName && source.Namespace == srcPod.Namespace {
			return true
		}
	}
	return false
}
