package envoy

import (
	"testing"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestHasDestinationEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		pass   bool
	}{
		{
			name: "no endpoints",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{},
			},
			pass: false,
		},
		{
			name: "no lb_endpoints",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{
					DynamicEndpointConfigs: []*adminv3.EndpointsConfigDump_DynamicEndpointConfig{
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{},
									},
								},
							}),
						},
					},
				},
			},
			pass: false,
		},
		{
			name: "one endpoint",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{
					DynamicEndpointConfigs: []*adminv3.EndpointsConfigDump_DynamicEndpointConfig{
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "1.2.3.4",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			pass: true,
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
			clusterChecker := HasDestinationEndpoints(configGetter)
			outcome := clusterChecker.Run()
			if test.pass {
				assert.NoError(outcome.GetError())
			} else {
				assert.Error(outcome.GetError())
			}
		})
	}
}

func TestHasSpecificEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		pod    *corev1.Pod
		pass   bool
	}{
		{
			name: "no endpoints in config",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{},
			},
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					PodIP: "1.2.3.4",
				},
			},
			pass: false,
		},
		{
			name: "pod IP matches no endpoint",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{
					DynamicEndpointConfigs: []*adminv3.EndpointsConfigDump_DynamicEndpointConfig{
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					PodIP: "1.2.3.4",
				},
			},
			pass: false,
		},
		{
			name: "pod IP matches only endpoint",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{
					DynamicEndpointConfigs: []*adminv3.EndpointsConfigDump_DynamicEndpointConfig{
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "1.2.3.4",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					PodIP: "1.2.3.4",
				},
			},
			pass: true,
		},
		{
			name: "pod IP matches one of several endpoints",
			config: &Config{
				Endpoints: adminv3.EndpointsConfigDump{
					DynamicEndpointConfigs: []*adminv3.EndpointsConfigDump_DynamicEndpointConfig{
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
										},
									},
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}),
						},
						{
							EndpointConfig: marshalClusterLoadAssignmentOrDie(&endpointv3.ClusterLoadAssignment{
								Endpoints: []*endpointv3.LocalityLbEndpoints{
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
										},
									},
									{
										LbEndpoints: []*endpointv3.LbEndpoint{
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "0.0.0.0",
																},
															},
														},
													},
												},
											},
											{
												HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
													Endpoint: &endpointv3.Endpoint{
														Address: &corev3.Address{
															Address: &corev3.Address_SocketAddress{
																SocketAddress: &corev3.SocketAddress{
																	Address: "1.2.3.4",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					PodIP: "1.2.3.4",
				},
			},
			pass: true,
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
			clusterChecker := HasSpecificEndpoint(configGetter, test.pod)
			outcome := clusterChecker.Run()
			if test.pass {
				assert.NoError(outcome.GetError())
			} else {
				assert.Error(outcome.GetError())
			}
		})
	}
}

func marshalClusterLoadAssignmentOrDie(cla *endpointv3.ClusterLoadAssignment) *any.Any {
	a, err := ptypes.MarshalAny(cla)
	if err != nil {
		panic(err)
	}
	return a
}
