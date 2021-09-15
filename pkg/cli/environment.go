package cli

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/openservicemesh/osm-health/pkg/common"
)

const (
	defaultOSMNamespace = "osm-system"
	osmNamespaceEnvVar  = "OSM_NAMESPACE"
)

// EnvSettings describes all CLI environment settings
type EnvSettings struct {
	namespace string
	config    *genericclioptions.ConfigFlags
}

// New relevant environment variables set and returns EnvSettings
func New() *EnvSettings {
	env := &EnvSettings{
		namespace: envOr(osmNamespaceEnvVar, defaultOSMNamespace),
	}

	// bind to kubernetes config flags
	env.config = &genericclioptions.ConfigFlags{
		Namespace: &env.namespace,
	}
	return env
}

func envOr(name, defaultVal string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return defaultVal
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.namespace, "osm-namespace", s.namespace, "namespace for osm control plane")
}

// RESTClientGetter gets the kubeconfig from EnvSettings
func (s *EnvSettings) RESTClientGetter() genericclioptions.RESTClientGetter {
	return s.config
}

// Namespace gets the namespace from the configuration
func (s *EnvSettings) Namespace() common.MeshNamespace {
	if ns, _, err := s.config.ToRawKubeConfigLoader().Namespace(); err == nil {
		return common.MeshNamespace(ns)
	}
	return "default"
}
