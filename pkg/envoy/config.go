package envoy

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// tempStruct implements ConfigGetter and is temporary until https://github.com/openservicemesh/osm-health/pull/14 is merged
// TODO(draychev): delete this
type tempStruct struct {
	*v1.Pod
}

// GetConfig implements ConfigGetter
func (mcg tempStruct) GetConfig() (*Config, error) {
	return &Config{}, nil
}

// GetObjectName implements ConfigGetter
func (mcg tempStruct) GetObjectName() string {
	return fmt.Sprintf("%s/%s", mcg.Pod.Namespace, mcg.Pod.Name)
}

// GetEnvoyConfigGetterForPod returns a ConfigGetter struct, which can fetch the Envoy config for the given pod.
func GetEnvoyConfigGetterForPod(pod *v1.Pod) (ConfigGetter, error) {
	// TODO(draychev): Get the config from the Envoy sidecar of the given pod.
	return tempStruct{Pod: pod}, nil
}
