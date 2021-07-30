package envoy

import (
	v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	"github.com/openservicemesh/osm-health/pkg/logger"
)

var log = logger.New("osm-health/envoy")

// ConfigGetter is an interface for getting Envoy config from Pods' sidecars.
type ConfigGetter interface {
	// GetConfig returns Envoy config.
	GetConfig() (*Config, error)

	// GetObjectName returns the name of the object (Pod) from which we fetch Envoy config.
	GetObjectName() string
}

// Config is Envoy config dump.
type Config struct {
	// Boostrap is an Envoy xDS proto.
	Boostrap v3.BootstrapConfigDump

	// Clusters is an Envoy xDS proto.
	Clusters v3.ClustersConfigDump

	// Endpoints is an Envoy xDS proto.
	Endpoints v3.EndpointsConfigDump

	// Listeners is an Envoy xDS proto.
	Listeners v3.ListenersConfigDump

	// Routes is an Envoy xDS proto.
	Routes v3.RoutesConfigDump
}
