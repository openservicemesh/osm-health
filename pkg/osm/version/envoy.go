package version

// EnvoyAdminPort is the admin port number of Envoys configured by the given version of the OSM Controller.
var EnvoyAdminPort = map[ControllerVersion]uint16{
	"v0.5":  15000,
	"v0.6":  15000,
	"v0.7":  15000,
	"v0.8":  15000,
	"v0.9":  15000,
	"v0.10": 15000,
	"v0.11": 15000,
}
