package envoy

import k8s "github.com/openservicemesh/osm-health/pkg/kubernetes"

// GetEnvoyConfig returns the Envoy config for the given pod.
func GetEnvoyConfig(pod k8s.Pod) (*Config, error) {
	// TODO(draychev): Get the config from the Envoy sidecar of the given pod.
	return nil, nil
}
