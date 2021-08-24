package envoy

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"

	"github.com/openservicemesh/osm-health/pkg/osm"
)

func TestEnvoyListenerChecker(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("v0.9")
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore.json"),
	}
	listenerChecker := HasInboundListener(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.Nil(outcome.GetError())
}

func TestEnvoyListenerCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("v0.9")
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	listenerChecker := HasOutboundListener(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal("envoy config is empty", outcome.GetError().Error())
}

func TestEnvoyListenerCheckerInvalidOSMVersion(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("no-such-version")
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	listenerChecker := HasOutboundListener(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal("osm controller version not recognized", outcome.GetError().Error())
}
