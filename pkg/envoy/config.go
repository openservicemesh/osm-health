package envoy

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/osm/version"
	osmCLI "github.com/openservicemesh/osm/pkg/cli"
)

// ConfigGetterStruct implements ConfigGetter interface.
type ConfigGetterStruct struct {
	*corev1.Pod
	version.ControllerVersion
}

// GetConfig implements ConfigGetter interface.
func (mcg ConfigGetterStruct) GetConfig() (*Config, error) {
	client, err := pod.GetKubeClient()
	if err != nil {
		return nil, err
	}

	config, err := pod.GetKubeConfig()
	if err != nil {
		return nil, err
	}

	namespace := mcg.Pod.Namespace
	podName := mcg.Pod.Name
	localPort := version.EnvoyAdminPort[mcg.ControllerVersion]
	query := "config_dump?include_eds"
	// This function becomes available in github.com/openservicemesh/osm at 9be251135819c360ce2b9cf77087c88ab1e3f54a
	configBytes, err := osmCLI.GetEnvoyProxyConfig(client, config, namespace, podName, localPort, query)
	if err != nil {
		return nil, err
	}

	return ParseEnvoyConfig(configBytes)
}

// GetObjectName implements ConfigGetter
func (mcg ConfigGetterStruct) GetObjectName() string {
	return fmt.Sprintf("%s/%s", mcg.Pod.Namespace, mcg.Pod.Name)
}

// GetEnvoyConfigGetterForPod returns a ConfigGetter struct, which can fetch the Envoy config for the given pod.
func GetEnvoyConfigGetterForPod(pod *corev1.Pod, osmVersion version.ControllerVersion) (ConfigGetter, error) {
	return ConfigGetterStruct{
		Pod:               pod,
		ControllerVersion: osmVersion,
	}, nil
}
