package access

import (
	"github.com/openservicemesh/osm-health/pkg/osm"
)

// isTrafficTargetRouteKindSupported checks whether an SMI TrafficTarget Route Kind is supported.
func isTrafficTargetRouteKindSupported(routeKind string, osmVersion osm.ControllerVersion) error {
	supportedRouteKinds, ok := osm.SupportedTrafficTargetRouteKinds[osmVersion]
	if !ok {
		return ErrorUnknownSupportForRouteKindUnknownOsmVersion
	}
	for _, supportedRouteKind := range supportedRouteKinds {
		if routeKind == string(supportedRouteKind) {
			return nil
		}
	}
	return ErrorUnsupportedRouteKind
}
