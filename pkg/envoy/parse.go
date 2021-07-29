package envoy

import (
	v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	// All of these are required for JSON to ConfigDump parsing to work
	_ "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/stream/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/wasm/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"

	"github.com/golang/protobuf/ptypes"
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
		log.Err(err).Msg("Error parsing JSON bytes")
		return nil, err
	}

	var err error
	var cfg Config

	err = ptypes.UnmarshalAny(configDump.Configs[0], &cfg.Boostrap)
	if err != nil {
		log.Err(err).Msg("Error parsing Bootstrap")
		return nil, err
	}

	err = ptypes.UnmarshalAny(configDump.Configs[1], &cfg.Clusters)
	if err != nil {
		log.Err(err).Msg("Error parsing Clusters")
		return nil, err
	}

	err = ptypes.UnmarshalAny(configDump.Configs[2], &cfg.Listeners)
	if err != nil {
		log.Err(err).Msg("Error parsing Listeners")
		return nil, err
	}

	err = ptypes.UnmarshalAny(configDump.Configs[4], &cfg.Routes)
	if err != nil {
		log.Err(err).Msg("Error parsing Routes")
		return nil, err
	}

	return &cfg, nil
}
