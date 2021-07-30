package envoy

import (
	"fmt"
	"os"
	"testing"

	tassert "github.com/stretchr/testify/assert"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

func TestEnvoyConfigParser(t *testing.T) {
	assert := tassert.New(t)
	sampleJsonfileName := "../../tests/sample-enovy-config-dump.json"
	sampleConfig, err := os.ReadFile(sampleJsonfileName)
	assert.Nil(err, fmt.Sprintf("Error opening file %s: %s", sampleJsonfileName, err))
	cfg, err := ParseEnvoyConfig(sampleConfig)
	assert.Nil(err, fmt.Sprintf("Error parsing Envoy config dump: %s", err))
	assert.NotNil(cfg, "Parsed Envoy config is empty.")

	// Bootstrap
	{
		assert.Equal(cfg.Boostrap.Bootstrap.Node.Id, "38cf2479-bfea-4c1e-a961-f8f8e2b2e8cb.sidecar.bookbuyer.bookbuyer.cluster.local")
	}

	{
		// Clusters
		assert.Len(cfg.Clusters.DynamicActiveClusters, 1)
		var actual envoy_config_cluster_v3.Cluster
		err := cfg.Clusters.DynamicActiveClusters[0].Cluster.UnmarshalTo(&actual)
		assert.Nil(err)
		assert.Equal(actual.Name, "bookstore/bookstore")
	}

	{
		// Listeners
		assert.Len(cfg.Listeners.DynamicListeners, 1)
		actual := cfg.Listeners.DynamicListeners[0]
		assert.Equal(actual.Name, "outbound-listener")
	}

	{
		// Routes
		assert.Len(cfg.Routes.DynamicRouteConfigs, 1)
		var actual envoy_config_route_v3.RouteConfiguration
		err := cfg.Routes.DynamicRouteConfigs[0].RouteConfig.UnmarshalTo(&actual)
		assert.Nil(err)
		assert.Equal(actual.Name, "rds-outbound")
	}
}
