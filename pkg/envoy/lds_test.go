package envoy

import (
	"os"
	"testing"

	tassert "github.com/stretchr/testify/assert"

	"github.com/openservicemesh/osm-health/pkg/osm"
)

type mockConfigGetter struct {
	getter func() (*Config, error)
}

func (mcg mockConfigGetter) GetConfig() (*Config, error) {
	return mcg.getter()
}

func (mcg mockConfigGetter) GetObjectName() string {
	return "namespace/podName"
}

func TestEnvoyListenerChecker(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("v0.9")
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			sampleConfig, err := os.ReadFile("../../tests/sample-enovy-config-dump.json")
			if err != nil {
				return nil, err
			}
			return ParseEnvoyConfig(sampleConfig)
		},
	}
	listenerChecker := HasListener(configGetter, osmVersion)
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
	listenerChecker := HasListener(configGetter, osmVersion)
	checkError := listenerChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("envoy config is empty", checkError.Error())
}

func TestEnvoyListenerCheckerInvalidOSMVersion(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("no-such-version")
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			sampleConfig, err := os.ReadFile("../../tests/sample-enovy-config-dump.json")
			if err != nil {
				return nil, err
			}
			return ParseEnvoyConfig(sampleConfig)
		},
	}
	listenerChecker := HasListener(configGetter, osmVersion)
	checkError := listenerChecker.Run()
	assert.NotNil(checkError)
	assert.Equal("osm controller version not recognized", checkError.Error())
}
