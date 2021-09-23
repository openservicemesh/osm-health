package version

// ControllerVersion is a string type alias for the OSM version.
type ControllerVersion string

// String() implements the stringer for ControllerVersion.
func (v ControllerVersion) String() string {
	return string(v)
}

// IngressVersion is a string type alias for the Ingress version supported.
type IngressVersion string

// TrafficTargetVersion is a string type alias for the SMI TrafficTarget version supported.
type TrafficTargetVersion string

// TrafficTargetRouteKind is a string type alias for the SMI TrafficTarget route kind supported.
type TrafficTargetRouteKind string

// TrafficSplitVersion is a string type alias for the SMI TrafficSplit version supported.
type TrafficSplitVersion string

// HTTPRouteVersion is a string type alias for the SMI HTTPRoute version supported.
type HTTPRouteVersion string

// Annotation is a string type alias for the osm annotation.
type Annotation string
