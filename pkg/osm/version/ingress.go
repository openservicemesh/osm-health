package version

// SupportedIngress maintains a mapping of OSM version to supported Ingress resource versions.
var SupportedIngress = map[ControllerVersion][]IngressVersion{
	"v0.5": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.5/pkg/ingress/client.go#L6
		"extensions/v1beta1",
	},
	"v0.6": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.6/pkg/ingress/client.go#L6
		"extensions/v1beta1",
	},
	"v0.7": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.7/pkg/ingress/client.go#L6
		"extensions/v1beta1",
	},
	"v0.8": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.8/pkg/ingress/client.go#L6
		"networking/v1beta1",
	},
	"v0.9": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.9/pkg/ingress/client.go#L6-L7
		"networking/v1",
		"networking/v1beta1",
	},
	"v0.10": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.10/pkg/ingress/client.go#L5-L6
		"networking/v1",
		"networking/v1beta1",
	},
	"v0.11": {
		// Source: https://github.com/openservicemesh/osm/blob/release-v0.11/pkg/ingress/client.go#L5-L6
		"networking/v1",
		"networking/v1beta1",
	},
}
