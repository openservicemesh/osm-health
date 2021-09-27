package envoy

import (
	"fmt"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/runner"
)

// Verify interface compliance
var _ runner.Runnable = (*DynamicWarmingCheck)(nil)

// DynamicWarmingCheck implements common.Runnable
type DynamicWarmingCheck struct {
	ConfigGetter
}

// Run implements common.Runnable
func (l DynamicWarmingCheck) Run() outcomes.Outcome {
	if l.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	envoyConfig, err := l.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	if len(envoyConfig.SecretsConfigDump.DynamicWarmingSecrets) > 0 {
		return outcomes.Fail{Error: ErrDynamicWarmingSecretsConfigDumpNotEmpty}
	}

	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (l DynamicWarmingCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (l DynamicWarmingCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (l DynamicWarmingCheck) Description() string {
	return fmt.Sprintf("Checking whether %s has dynamic warming issues", l.ConfigGetter.GetObjectName())
}

// NewDynamicWarmingCheck creates a DynamicWarmingCheck which checks whether the given Pod's envoy has dynamic warming issues.
func NewDynamicWarmingCheck(configGetter ConfigGetter) DynamicWarmingCheck {
	return DynamicWarmingCheck{
		ConfigGetter: configGetter,
	}
}
