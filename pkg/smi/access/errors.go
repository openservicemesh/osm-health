package access

import "errors"

var (
	// ErrorUnknownSupportForRouteKindUnknownOsmVersion is the error for unknown osm versions when checking for supported traffic target route kinds.
	ErrorUnknownSupportForRouteKindUnknownOsmVersion = errors.New("unknown osm version: no info on supported traffic target route kinds for specified osm version")

	// ErrorUnsupportedRouteKind is the error if the osm version does not support the TrafficTarget route kind.
	ErrorUnsupportedRouteKind = errors.New("unsupported traffic target route kind")
)
