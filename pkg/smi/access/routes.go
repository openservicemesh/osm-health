package access

import (
	"github.com/openservicemesh/osm-health/pkg/osm/version"
)

// isTrafficTargetRouteKindSupported checks whether an SMI TrafficTarget Route Kind is supported.
func isTrafficTargetRouteKindSupported(routeKind string, osmVersion version.ControllerVersion) error {
	supportedRouteKinds, ok := version.SupportedTrafficTargetRouteKinds[osmVersion]
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
