package envoy

import (
	"os"
	"testing"

	tassert "github.com/stretchr/testify/assert"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

func configFromFileOrFail(t *testing.T, filename string) *Config {
	t.Helper()
	sampleConfig, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Error opening %s: %v", filename, err)
	}
	cfg, err := ParseEnvoyConfig(sampleConfig)
	if err != nil {
		t.Fatal("Error parsing Envoy config dump:", err)
	}
	if cfg == nil {
		t.Fatal("Parsed Envoy config is empty")
	}
	return cfg
}

func TestEnvoyConfigParserBookbuyer(t *testing.T) {
	cfg := configFromFileOrFail(t, "../../tests/sample-enovy-config-dump-bookbuyer.json")

	assert := tassert.New(t)
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

func TestEnvoyConfigParserBookstore(t *testing.T) {
	cfg := configFromFileOrFail(t, "../../tests/sample-enovy-config-dump-bookstore.json")

	assert := tassert.New(t)
	// Bootstrap
	{
		assert.Equal("b2d941c7-484a-4cd4-ad65-76e41b79e48a.sidecar.bookstore-v1.bookstore.cluster.local", cfg.Boostrap.Bootstrap.Node.Id)
	}

	{
		// Clusters
		assert.Len(cfg.Clusters.DynamicActiveClusters, 3)
		var actual envoy_config_cluster_v3.Cluster
		err := cfg.Clusters.DynamicActiveClusters[0].Cluster.UnmarshalTo(&actual)
		assert.Nil(err)
		assert.Equal("bookstore/bookstore-v1-local", actual.Name)
	}

	{
		// Listeners
		assert.Len(cfg.Listeners.DynamicListeners, 3)
		actual := cfg.Listeners.DynamicListeners[0]
		assert.Equal("outbound-listener", actual.Name)
	}

	{
		// Routes
		assert.Len(cfg.Routes.DynamicRouteConfigs, 2)
		var actual envoy_config_route_v3.RouteConfiguration
		err := cfg.Routes.DynamicRouteConfigs[0].RouteConfig.UnmarshalTo(&actual)
		assert.Nil(err)
		assert.Equal("rds-outbound", actual.Name)
	}
}
