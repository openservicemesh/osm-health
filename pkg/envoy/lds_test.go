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
<<<<<<< HEAD
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore.json"),
=======
		getter: createConfigGetterFuncFromConfigFile("../../tests/sample-envoy-config-dump-bookstore.json"),
>>>>>>> feat(connectivity): Check rds cfg domains
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
<<<<<<< HEAD
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
=======
		getter: createConfigGetterFuncFromConfigFile("../../tests/sample-envoy-config-dump-bookbuyer.json"),
>>>>>>> feat(connectivity): Check rds cfg domains
	}
	listenerChecker := HasOutboundListener(configGetter, osmVersion)
	checkError := listenerChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("osm controller version not recognized", checkError.Error())
}
