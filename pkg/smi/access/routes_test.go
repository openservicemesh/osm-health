package access

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"

	"github.com/openservicemesh/osm-health/pkg/osm/version"
	"github.com/openservicemesh/osm-health/pkg/smi"
)

func TestIsTrafficTargetRouteKindSupported(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		routeKind string
		expErr    error
	}{
		{
			name:      "known osm version and supported route kind - HTTPRouteGroupKind",
			version:   "v0.9",
			routeKind: smi.HTTPRouteGroupKind,
			expErr:    nil,
		},
		{
			name:      "known osm version and supported route kind - TCPRouteKind",
			version:   "v0.9",
			routeKind: smi.TCPRouteKind,
			expErr:    nil,
		},
		{
			name:      "known osm version and unsupported route kind",
			version:   "v0.9",
			routeKind: "some-other-route-kind",
			expErr:    ErrorUnsupportedRouteKind,
		},
		{
			name:      "unknown osm version",
			version:   "v-abc-def-xxx",
			routeKind: "some-other-route-kind",
			expErr:    ErrorUnknownSupportForRouteKindUnknownOsmVersion,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			err := isTrafficTargetRouteKindSupported(test.routeKind, version.ControllerVersion(test.version))
			assert.Equal(test.expErr, err)
		})
	}
}
