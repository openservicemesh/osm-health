package osm

import "github.com/openservicemesh/osm-health/pkg/smi"

// SupportedTrafficTarget is a map of OSM Controller Version to supported
var SupportedTrafficTarget = map[ControllerVersion]TrafficTargetVersion{
	// Source: https://github.com/openservicemesh/osm/blob/release-v0.5/pkg/smi/client.go#L8
	"v0.5": "v1alpha2",

	// Source: https://github.com/openservicemesh/osm/blob/release-v0.6/pkg/smi/client.go#L8
	"v0.6": "v1alpha2",

	// Source: https://github.com/openservicemesh/osm/blob/release-v0.7/pkg/smi/client.go#L8
	"v0.7": "v1alpha3",

	// Source: https://github.com/openservicemesh/osm/blob/release-v0.8/pkg/smi/client.go#L8
	"v0.8": "v1alpha3",

	// Source: https://github.com/openservicemesh/osm/blob/release-v0.9/pkg/smi/client.go#L10
	"v0.9": "v1alpha3",
}

// SupportedTrafficTargetRouteKinds is a map of OSM Controller Version to supported SMI TrafficTarget Route Kinds.
var SupportedTrafficTargetRouteKinds = map[ControllerVersion][]TrafficTargetRouteKind{
	// Source:
	// https://github.com/openservicemesh/osm/blob/release-v0.5/pkg/smi/client.go#L67
	// https://github.com/openservicemesh/osm/blob/release-v0.5/pkg/smi/client.go#L68
	"v0.5": {
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
	},

	// Sources:
	// https://github.com/openservicemesh/osm/blob/release-v0.6/pkg/smi/client.go#L68
	// https://github.com/openservicemesh/osm/blob/release-v0.6/pkg/smi/client.go#L69
	"v0.6": {
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
	},

	// Sources:
	// https://github.com/openservicemesh/osm/blob/release-v0.7/pkg/smi/client.go#L68
	// https://github.com/openservicemesh/osm/blob/release-v0.7/pkg/smi/client.go#L69
	"v0.7": {
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
	},

	// Sources:
	// https://github.com/openservicemesh/osm/blob/release-v0.8/pkg/smi/client.go#L58
	// https://github.com/openservicemesh/osm/blob/release-v0.8/pkg/smi/client.go#L59
	"v0.8": {
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
	},

	// Sources:
	// https://github.com/openservicemesh/osm/blob/release-v0.9/pkg/smi/client.go#L60
	// https://github.com/openservicemesh/osm/blob/release-v0.9/pkg/smi/client.go#L61
	"v0.9": {
		smi.HTTPRouteGroupKind,
		smi.TCPRouteKind,
	},
}

// SupportedTrafficSplit is the mapping of OSM Controller version to supported SMI TrafficSplit version.
var SupportedTrafficSplit = map[ControllerVersion]TrafficSplitVersion{
	"v0.5": "v1alpha2",
	"v0.6": "v1alpha2",
	"v0.7": "v1alpha2",
	"v0.8": "v1alpha2",
	"v0.9": "v1alpha2",
}

// SupportedHTTPRouteVersion is a mapping of OSM Controller version to supported HTTP Route Group version.
var SupportedHTTPRouteVersion = map[ControllerVersion][]HTTPRouteVersion{
	"v0.5": {
		"v1alpha3",
	},
	"v0.6": {
		"v1alpha3",
	},
	"v0.7": {
		"v1alpha4",
	},
	"v0.8": {
		"v1alpha4",
	},
	"v0.9": {
		"v1alpha4",
	},
}
