package envoy

import (
	"testing"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	tassert "github.com/stretchr/testify/assert"
)

func TestDynamicWarmingCheck(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr error
	}{
		{
			name:        "no config",
			config:      nil,
			expectedErr: ErrEnvoyConfigEmpty,
		},
		{
			name: "no dynamic warming secrets in secrets config dump",
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{},
			},
			expectedErr: nil,
		},
		{
			name: "nil dynamic warming secrets in secrets config dump",
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicWarmingSecrets: nil,
				},
			},
			expectedErr: nil,
		},
		{
			name: "empty length dynamic warming secrets in secrets config dump",
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicWarmingSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{},
				},
			},
			expectedErr: nil,
		},
		{
			name: "dynamic warming secrets present in secrets config dump",
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicWarmingSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: "dynamic-warming-secret",
						},
					},
				},
			},
			expectedErr: ErrDynamicWarmingSecretsConfigDumpNotEmpty,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			configGetter := mockConfigGetter{
				getter: func() (*Config, error) {
					return test.config, nil
				},
			}
			dynamicWarmingChecker := NewDynamicWarmingCheck(configGetter)
			outcome := dynamicWarmingChecker.Run()
			if test.expectedErr == nil {
				assert.Nil(outcome.GetError())
			} else {
				assert.Equal(test.expectedErr.Error(), outcome.GetError().Error())
			}
		})
	}
}
