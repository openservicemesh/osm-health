package envoy

import v1 "k8s.io/api/core/v1"

// GetEnvoyConfig returns the Envoy config for the given pod.
func GetEnvoyConfig(pod *v1.Pod) (*Config, error) {
	// TODO(draychev): Get the config from the Envoy sidecar of the given pod.
	return nil, nil
}
