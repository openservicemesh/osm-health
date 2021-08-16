package envoy

import (
	"os"
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

func createConfigGetterFunc(configFilePath string) func() (*Config, error) {
	return func() (*Config, error) {
		sampleConfig, err := os.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}
		return ParseEnvoyConfig(sampleConfig)
	}
}
