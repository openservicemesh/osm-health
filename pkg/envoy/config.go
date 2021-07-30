package envoy

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm-health/pkg/osm"
	osmCLI "github.com/openservicemesh/osm/pkg/cli"
)

// ConfigGetterStruct implements ConfigGetter interface.
type ConfigGetterStruct struct {
	*v1.Pod
	osm.ControllerVersion
}

// GetConfig implements ConfigGetter interface.
func (mcg ConfigGetterStruct) GetConfig() (*Config, error) {
	client, err := kuberneteshelper.GetKubeClient()
	if err != nil {
		return nil, err
	}

	config, err := kuberneteshelper.GetKubeConfig()
	if err != nil {
		return nil, err
	}

	namespace := mcg.Pod.Namespace
	podName := mcg.Pod.Name
	localPort := osm.EnvoyAdminPort[mcg.ControllerVersion]
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
func GetEnvoyConfigGetterForPod(pod *v1.Pod, osmVersion osm.ControllerVersion) (ConfigGetter, error) {
	return ConfigGetterStruct{
		Pod:               pod,
		ControllerVersion: osmVersion,
	}, nil
}
