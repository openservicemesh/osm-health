package osm

// SupportedTrafficTarget is a map of OSM Controller Version to supported
var SupportedTrafficTarget = map[ControllerVersion][]TrafficTargetVersion{
	"v0.5": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.5/pkg/smi/client.go#L8
		"v1alpha2",
	},
	"v0.6": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.6/pkg/smi/client.go#L8
		"v1alpha2",
	},
	"v0.7": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.7/pkg/smi/client.go#L8
		"v1alpha3",
	},
	"v0.8": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.8/pkg/smi/client.go#L8
		"v1alpha3",
	},
	"v0.9": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.9/pkg/smi/client.go#L10
		"v1alpha3",
	},
}
