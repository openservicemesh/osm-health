package envoy

import (
	v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	"github.com/openservicemesh/osm-health/pkg/logger"
)

var log = logger.New("osm-health/connectivity")

// Config is Envoy config dump.
type Config struct {
	Boostrap  v3.BootstrapConfigDump
	Clusters  v3.ClustersConfigDump
	Endpoints v3.EndpointsConfigDump
	Listeners v3.ListenersConfigDump
	Routes    v3.RoutesConfigDump
}
