package connectivity

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodToPod tests the connectivity between a source and destination pods.
func PodToPod(srcPod *v1.Pod, destPod *v1.Pod, clientSet kubernetes.Interface) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s/%s", srcPod.Namespace, srcPod.Name, destPod.Namespace, destPod.Name)

	// TODO: actually test connectivity
	srcPodLabels := srcPod.ObjectMeta.GetLabels()
	for label,val := range srcPodLabels {
		fmt.Printf("label: %s, val: %s\n", label,val)
	}
	srcPodAnnotations := srcPod.ObjectMeta.GetAnnotations()
	for annotation,val := range srcPodAnnotations {
		fmt.Printf("annotation: %s, val: %s\n", annotation,val)
	}
	inSameMesh, err := arePodsInSameMesh(srcPod, destPod, clientSet)
	if err != nil {
		log.Err(err).Msg("Error getting list of Pods")
	}
	fmt.Printf("Pods are in same mesh: %t\n", inSameMesh)

	return common.Result{
		SMIPolicy: common.SMIPolicy{
			HasPolicy:                  false,
			ValidPolicy:                false,
			SourceToDestinationAllowed: false,
		},
		Successful: false,
	}
}

func arePodsInSameMesh(srcPod *v1.Pod, destPod *v1.Pod, clientSet kubernetes.Interface) (bool, error) {
	srcPodNamespace, err := clientSet.CoreV1().Namespaces().Get(context.Background(), srcPod.ObjectMeta.GetNamespace(), v12.GetOptions{})
	if err != nil {
		log.Err(err).Msg("Error getting source pod's namespace")
		return false, errors.New("error getting namespace")
	}
	destPodNamespace, err := clientSet.CoreV1().Namespaces().Get(context.Background(), destPod.ObjectMeta.GetNamespace(), v12.GetOptions{})
	if err != nil {
		log.Err(err).Msg("Error getting destination pod's namespace")
		return false, errors.New("error getting namespace")
	}
	srcPodMeshName := srcPodNamespace.ObjectMeta.GetLabels()["openservicemesh.io/monitored-by"]
	if destPodNamespace.ObjectMeta.GetLabels()["openservicemesh.io/monitored-by"] == srcPodMeshName {
		return true, nil
	}
	return false, nil
}
