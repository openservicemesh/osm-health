package kubernetesHelper

import (
	"context"
	"errors"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// PodFromString validates the name of the Pod
func PodFromString(namespacedPod string) (*v1.Pod, error) {
	podChunks := strings.Split(namespacedPod, "/")
	if len(podChunks) != 2 {
		log.Fatal().Msgf("Invalid Pod name %s; This is expected to be in the format: namespace/name", namespacedPod)
		return nil, errors.New("invalid Pod name")
	}

	if len(os.Getenv("KUBECONFIG")) <= 0 {
		log.Fatal().Msgf("Point us to the Kubernetes Config via the KUBECONFIG. export KUBECONFIG=~/.kube/config maybe?")
	}

	// Initialize kube config and client
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating kube configs using in-cluster config")
	}

	namespace := podChunks[0]
	podName := podChunks[1]

	log.Trace().Msgf("Looking for Pod with Name=%s in Namespace=%s", podName, namespace)

	kubeClient := kubernetes.NewForConfigOrDie(kubeConfig)

	podList, err := kubeClient.CoreV1().Pods(namespace).List(context.Background(), v12.ListOptions{})
	if err != nil {
		log.Err(err).Msg("Error getting list of Pods")
		return nil, errors.New("error getting pods")
	}

	log.Trace().Msgf("Looking for pod %s", namespacedPod)
	for _, pod := range podList.Items {
		log.Trace().Msgf("Could this be it: %s/%s", pod.Namespace, pod.Name)
		if pod.Namespace == namespace && pod.Name == podName {
			log.Trace().Msgf("Found Pod %s: %+v", namespacedPod, pod)
			return &pod, nil
		}
	}

	log.Error().Msgf("Did not find Pod %s", namespacedPod)
	return nil, errors.New("no pod found")
}
