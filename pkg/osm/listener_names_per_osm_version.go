package osm

// OutboundListenerNames is the name of the Envoy listener expected to be created by a certain version of the OSM Controller.
var OutboundListenerNames = map[ControllerVersion]string{
	"v0.5": "outbound-listener",
	"v0.6": "outbound-listener",
	"v0.7": "outbound-listener",
	"v0.8": "outbound-listener",
	"v0.9": "outbound-listener",
}

// InboundListenerNames is the name of the Envoy listener expected to be created by a certain version of the OSM Controller.
var InboundListenerNames = map[ControllerVersion]string{
	"v0.5": "inbound-listener",
	"v0.6": "inbound-listener",
	"v0.7": "inbound-listener",
	"v0.8": "inbound-listener",
	"v0.9": "inbound-listener",
}
