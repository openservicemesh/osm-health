package kuberneteshelper

import (
	"context"
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned"
	"github.com/openservicemesh/osm/pkg/signals"
)

// PodFromString validates the name of the Pod
func PodFromString(namespacedPod string) (*v1.Pod, error) {
	podChunks := strings.Split(namespacedPod, "/")
	if len(podChunks) != 2 {
		log.Fatal().Msgf("Invalid Pod name %s; This is expected to be in the format: namespace/name", namespacedPod)
		return nil, errors.New("invalid Pod name")
	}

	namespace := podChunks[0]
	podName := podChunks[1]

	log.Trace().Msgf("Looking for Pod with Name=%s in namespace=%s", podName, namespace)

	kubeClient, err := GetKubeClient()
	if err != nil {
		log.Err(err).Msgf("Error getting Kubernetes client")
		return nil, err
	}

	podList, err := kubeClient.CoreV1().Pods(namespace).List(context.Background(), v12.ListOptions{})
	if err != nil {
		log.Err(err).Msg("Error getting list of Pods")
		return nil, errors.New("error getting pods")
	}

	log.Trace().Msgf("Looking for pod %s", namespacedPod)
	for _, pod := range podList.Items {
		if pod.Namespace == namespace && pod.Name == podName {
			log.Trace().Msgf("Found Pod %s/%s", pod.Namespace, pod.Name)
			return &pod, nil
		}
	}

	log.Error().Msgf("Did not find Pod %s", namespacedPod)
	return nil, errors.New("no pod found")
}

// GetKubeConfig returns the kubeconfig
func GetKubeConfig() (*restclient.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}).ClientConfig()
}

// GetKubeClient returns a Kubernetes clientset.
func GetKubeClient() (kubernetes.Interface, error) {
	kubeConfig, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfigOrDie(kubeConfig), nil
}

// GetOsmConfigurator returns a new OSM configurator
func GetOsmConfigurator(osmNamespace common.MeshNamespace) configurator.Configurator {
	stop := signals.RegisterExitHandlers()
	kubeConfig, err := GetKubeConfig()
	if err != nil {
		log.Err(err).Msg("Error getting kubeconfig")
	}
	cfg := configurator.NewConfigurator(versioned.NewForConfigOrDie(kubeConfig), stop, osmNamespace.String(), constants.OSMMeshConfig)
	return cfg
}
