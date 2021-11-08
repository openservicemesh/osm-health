package envoy

import (
	v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	// All of these are required for JSON to ConfigDump parsing to work
	_ "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/stream/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/health_check/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/wasm/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/rbac/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"

	"google.golang.org/protobuf/encoding/protojson"
)

// ParseEnvoyConfig parses Envoy config_dump
func ParseEnvoyConfig(jsonBytes []byte) (*Config, error) {
	var configDump v3.ConfigDump
	unmarshal := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	if err := unmarshal.Unmarshal(jsonBytes, &configDump); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON bytes")
		return nil, err
	}

	var cfg Config

	for idx, config := range configDump.Configs {
		switch config.TypeUrl {
		case "type.googleapis.com/envoy.admin.v3.BootstrapConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.Boostrap); err != nil {
				log.Error().Err(err).Msg("Error parsing Bootstrap")
				return nil, err
			}

		case "type.googleapis.com/envoy.admin.v3.ClustersConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.Clusters); err != nil {
				log.Error().Err(err).Msg("Error parsing Clusters")
				return nil, err
			}

		case "type.googleapis.com/envoy.admin.v3.EndpointsConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.Endpoints); err != nil {
				log.Error().Err(err).Msg("Error parsing Endpoints")
				return nil, err
			}

		case "type.googleapis.com/envoy.admin.v3.ListenersConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.Listeners); err != nil {
				log.Error().Err(err).Msg("Error parsing Listeners")
				return nil, err
			}

		case "type.googleapis.com/envoy.admin.v3.RoutesConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.Routes); err != nil {
				log.Error().Err(err).Msg("Error parsing Listeners")
				return nil, err
			}
		case "type.googleapis.com/envoy.admin.v3.ScopedRoutesConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.ScopedRoutesConfigDump); err != nil {
				log.Error().Err(err).Msg("Error parsing ScopedRoutesConfigDump")
				return nil, err
			}

		case "type.googleapis.com/envoy.admin.v3.SecretsConfigDump":
			if err := configDump.Configs[idx].UnmarshalTo(&cfg.SecretsConfigDump); err != nil {
				log.Error().Err(err).Msg("Error parsing SecretsConfigDump")
				return nil, err
			}

		default:
			log.Error().Msgf("Unrecognized TypeUrl %s", config.TypeUrl)
		}
	}

	return &cfg, nil
}
