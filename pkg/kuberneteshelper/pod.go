package kuberneteshelper

import (
	"context"
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
func PodFromString(namespacedPod string) (*corev1.Pod, error) {
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
		log.Error().Err(err).Msgf("Error getting Kubernetes client")
		return nil, err
	}

	podList, err := kubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error getting list of Pods")
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
		log.Error().Err(err).Msg("Error getting kubeconfig")
	}
	cfg := configurator.NewConfigurator(versioned.NewForConfigOrDie(kubeConfig), stop, osmNamespace.String(), constants.OSMMeshConfig)
	return cfg
}

// GetMatchingServices returns a list of Kubernetes services in the namespace that match the pod's label
func GetMatchingServices(kubeClient kubernetes.Interface, podLabels map[string]string, namespace string) ([]*corev1.Service, error) {
	var serviceList []*corev1.Service
	svcList, err := kubeClient.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, svc := range svcList.Items {
		svc := svc
		svcRawSelector := svc.Spec.Selector
		// service has no selectors, we do not need to match against the pod label
		if len(svcRawSelector) == 0 {
			continue
		}
		selector := labels.SelectorFromSet(svcRawSelector)
		if selector.Matches(labels.Set(podLabels)) {
			serviceList = append(serviceList, &svc)
		}
	}

	return serviceList, nil
}
