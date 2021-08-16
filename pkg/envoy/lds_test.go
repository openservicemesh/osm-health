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
	checkError := listenerChecker.Run()
	assert.Nil(checkError)
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
	checkError := listenerChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyListenerCheckerInvalidOSMVersion(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("no-such-version")
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	listenerChecker := HasOutboundListener(configGetter, osmVersion)
	checkError := listenerChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("osm controller version not recognized", checkError.Error())
}
