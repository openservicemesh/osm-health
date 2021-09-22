package envoy

import (
	"testing"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha2"
	v1a2 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha2"
	"github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha3"
	v1a3 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/access/v1alpha3"
	fakeAccess "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned/fake"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openservicemesh/osm-health/pkg/osm"
	"github.com/openservicemesh/osm-health/pkg/smi"
)

func TestEnvoyListenerChecker(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("v0.9")
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookstore.json"),
	}
	listenerChecker := NewInboundListenerCheck(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.Nil(outcome.GetError())
}

func TestEnvoyListenerCheckerEmptyConfig(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("v0.9")
	configGetter := mockConfigGetter{
		getter: func() (*Config, error) {
			return nil, nil
		},
	}
	listenerChecker := NewOutboundListenerCheck(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal("envoy config is empty", outcome.GetError().Error())
}

func TestEnvoyListenerCheckerInvalidOSMVersion(t *testing.T) {
	assert := tassert.New(t)
	osmVersion := osm.ControllerVersion("no-such-version")
	configGetter := mockConfigGetter{
		getter: createConfigGetterFunc("../../tests/sample-envoy-config-dump-bookbuyer.json"),
	}
	listenerChecker := NewOutboundListenerCheck(configGetter, osmVersion)
	outcome := listenerChecker.Run()
	assert.NotNil(outcome.GetError())
	assert.Equal("osm controller version not recognized", outcome.GetError().Error())
}

func TestGetPossibleInboundFilterChainNames(t *testing.T) {
	tests := []struct {
		name                string
		ruleTypes           map[string]struct{}
		svcs                []*corev1.Service
		expFilterChainNames map[string]bool
		expErr              error
	}{
		{
			name: "no services",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
			},
			svcs:                nil,
			expFilterChainNames: map[string]bool{},
			expErr:              nil,
		},
		{
			name:      "no rules",
			ruleTypes: nil,
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{},
			expErr:              nil,
		},
		{
			name: "multiple services with one rule type",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{
				"inbound-mesh-http-filter-chain:mynamespace/myservice":  false,
				"inbound-mesh-http-filter-chain:mynamespace/myservice2": false,
			},
			expErr: nil,
		},
		{
			name: "multiple services with HTTP and TCP rules",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
				smi.TCPRouteKind:       {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{
				"inbound-mesh-http-filter-chain:mynamespace/myservice":  false,
				"inbound-mesh-tcp-filter-chain:mynamespace/myservice":   false,
				"inbound-mesh-http-filter-chain:mynamespace/myservice2": false,
				"inbound-mesh-tcp-filter-chain:mynamespace/myservice2":  false,
			},
			expErr: nil,
		},
		{
			name: "invalid rule kind",
			ruleTypes: map[string]struct{}{
				"unknownrule": {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: nil,
			expErr:              smi.ErrInvalidRuleKind,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			possibleFilterChainNames, err := getPossibleInboundFilterChainNames(test.svcs, test.ruleTypes)

			assert.Equal(test.expFilterChainNames, possibleFilterChainNames)
			assert.Equal(test.expErr, err)
		})
	}
}

func TestGetPossibleOutboundFilterChainNames(t *testing.T) {
	tests := []struct {
		name                string
		ruleTypes           map[string]struct{}
		svcs                []*corev1.Service
		expFilterChainNames map[string]bool
		expErr              error
	}{
		{
			name: "no services",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
			},
			svcs:                nil,
			expFilterChainNames: map[string]bool{},
			expErr:              nil,
		},
		{
			name:      "no rules",
			ruleTypes: nil,
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{},
			expErr:              nil,
		},
		{
			name: "multiple services with one rule type",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{
				"outbound-mesh-http-filter-chain:mynamespace/myservice":  false,
				"outbound-mesh-http-filter-chain:mynamespace/myservice2": false,
			},
			expErr: nil,
		},
		{
			name: "multiple services with HTTP and TCP rules",
			ruleTypes: map[string]struct{}{
				smi.HTTPRouteGroupKind: {},
				smi.TCPRouteKind:       {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: map[string]bool{
				"outbound-mesh-http-filter-chain:mynamespace/myservice":  false,
				"outbound-mesh-tcp-filter-chain:mynamespace/myservice":   false,
				"outbound-mesh-http-filter-chain:mynamespace/myservice2": false,
				"outbound-mesh-tcp-filter-chain:mynamespace/myservice2":  false,
			},
			expErr: nil,
		},
		{
			name: "invalid rule kind",
			ruleTypes: map[string]struct{}{
				"unknownrule": {},
			},
			svcs: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myservice2",
						Namespace: "mynamespace",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"mykey": "myval",
						},
					},
				},
			},
			expFilterChainNames: nil,
			expErr:              smi.ErrInvalidRuleKind,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			possibleFilterChainNames, err := getPossibleOutboundFilterChainNames(test.svcs, test.ruleTypes)

			assert.Equal(test.expFilterChainNames, possibleFilterChainNames)
			assert.Equal(test.expErr, err)
		})
	}
}

func TestGetRuleTypesFromMatchingTrafficTargetsV1alpha2(t *testing.T) {
	simpleSrcPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "srcpod",
			Namespace: "srcnamespace",
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "srcserviceaccount",
		},
	}

	simpleDstPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dstpod",
			Namespace: "dstnamespace",
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "dstserviceaccount",
		},
	}

	tests := []struct {
		name           string
		srcPod         *corev1.Pod
		dstPod         *corev1.Pod
		trafficTargets []*v1a2.TrafficTarget
		expRuleTypes   map[string]struct{}
		expErr         bool
	}{
		{
			name:         "no traffic targets",
			srcPod:       simpleSrcPod,
			dstPod:       simpleDstPod,
			expRuleTypes: map[string]struct{}{},
			trafficTargets: []*v1alpha2.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha2.TrafficTargetSpec{
						Destination: v1alpha2.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "notdstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha2.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha2.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "notsrcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:         "no matching traffic targets",
			srcPod:       simpleSrcPod,
			dstPod:       simpleDstPod,
			expRuleTypes: map[string]struct{}{},
			trafficTargets: []*v1alpha2.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha2.TrafficTargetSpec{
						Destination: v1alpha2.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "notdstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha2.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha2.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "notsrcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:   "one matching traffic targets",
			srcPod: simpleSrcPod,
			dstPod: simpleDstPod,
			expRuleTypes: map[string]struct{}{
				"HTTPRouteGroup": {},
			},
			trafficTargets: []*v1alpha2.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha2.TrafficTargetSpec{
						Destination: v1alpha2.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "dstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha2.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha2.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "srcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:   "one matching traffic targets and multiple rules",
			srcPod: simpleSrcPod,
			dstPod: simpleDstPod,
			expRuleTypes: map[string]struct{}{
				"HTTPRouteGroup": {},
				"TCPRoute":       {},
			},
			trafficTargets: []*v1alpha2.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha2.TrafficTargetSpec{
						Destination: v1alpha2.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "dstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha2.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule2",
							},
							{
								Kind: "TCPRoute",
								Name: "tcprouterule",
							},
						},
						Sources: []v1alpha2.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "srcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			objs := make([]runtime.Object, len(test.trafficTargets))
			for i := range test.trafficTargets {
				objs[i] = test.trafficTargets[i]
			}
			fakeAccessClient := fakeAccess.NewSimpleClientset(objs...)
			ruleTypes, err := getRuleTypesFromMatchingTrafficTargetsV1alpha2(test.srcPod, test.dstPod, fakeAccessClient)

			assert.Equal(test.expErr, err != nil)
			assert.Equal(test.expRuleTypes, ruleTypes)
		})
	}
}

func TestGetRuleTypesFromMatchingTrafficTargetsV1alpha3(t *testing.T) {
	simpleSrcPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "srcpod",
			Namespace: "srcnamespace",
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "srcserviceaccount",
		},
	}

	simpleDstPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dstpod",
			Namespace: "dstnamespace",
			Labels: map[string]string{
				"mykey": "myval",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "dstserviceaccount",
		},
	}

	tests := []struct {
		name           string
		srcPod         *corev1.Pod
		dstPod         *corev1.Pod
		trafficTargets []*v1a3.TrafficTarget
		expRuleTypes   map[string]struct{}
		expErr         bool
	}{
		{
			name:         "no traffic targets",
			srcPod:       simpleSrcPod,
			dstPod:       simpleDstPod,
			expRuleTypes: map[string]struct{}{},
			trafficTargets: []*v1alpha3.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha3.TrafficTargetSpec{
						Destination: v1alpha3.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "notdstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha3.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha3.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "notsrcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:         "no matching traffic targets",
			srcPod:       simpleSrcPod,
			dstPod:       simpleDstPod,
			expRuleTypes: map[string]struct{}{},
			trafficTargets: []*v1alpha3.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha3.TrafficTargetSpec{
						Destination: v1alpha3.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "notdstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha3.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha3.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "notsrcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:   "one matching traffic targets",
			srcPod: simpleSrcPod,
			dstPod: simpleDstPod,
			expRuleTypes: map[string]struct{}{
				"HTTPRouteGroup": {},
			},
			trafficTargets: []*v1alpha3.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha3.TrafficTargetSpec{
						Destination: v1alpha3.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "dstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha3.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
						},
						Sources: []v1alpha3.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "srcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name:   "one matching traffic targets and multiple rules",
			srcPod: simpleSrcPod,
			dstPod: simpleDstPod,
			expRuleTypes: map[string]struct{}{
				"HTTPRouteGroup": {},
				"TCPRoute":       {},
			},
			trafficTargets: []*v1alpha3.TrafficTarget{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "srcdsttraffictarget",
						Namespace: "dstnamespace",
					},
					Spec: v1alpha3.TrafficTargetSpec{
						Destination: v1alpha3.IdentityBindingSubject{
							Kind:      "ServiceAccount",
							Name:      "dstserviceaccount",
							Namespace: "dstnamespace",
						},
						Rules: []v1alpha3.TrafficTargetRule{
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule",
							},
							{
								Kind: "HTTPRouteGroup",
								Name: "httprouterule2",
							},
							{
								Kind: "TCPRoute",
								Name: "tcprouterule",
							},
						},
						Sources: []v1alpha3.IdentityBindingSubject{
							{
								Kind:      "ServiceAccount",
								Name:      "srcserviceaccount",
								Namespace: "srcnamespace",
							},
						},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			objs := make([]runtime.Object, len(test.trafficTargets))
			for i := range test.trafficTargets {
				objs[i] = test.trafficTargets[i]
			}
			fakeAccessClient := fakeAccess.NewSimpleClientset(objs...)
			ruleTypes, err := getRuleTypesFromMatchingTrafficTargetsV1alpha3(test.srcPod, test.dstPod, fakeAccessClient)

			assert.Equal(test.expErr, err != nil)
			assert.Equal(test.expRuleTypes, ruleTypes)
		})
	}
}

func TestFindMatchingFilterChainNames(t *testing.T) {
	tests := []struct {
		name                     string
		config                   *Config
		expectedListenerName     string
		possibleFilterChainNames map[string]bool
		expErr                   error
	}{
		{
			name: "listener name not found",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "not-outbound-listener",
							ActiveState: &adminv3.ListenersConfigDump_DynamicListenerState{
								Listener: marshalListenerOrDie(&listenerv3.Listener{
									Name: "not-outbound-listener",
									FilterChains: []*listenerv3.FilterChain{
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v1",
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedListenerName:     "outbound-listener",
			possibleFilterChainNames: map[string]bool{"outbound-mesh-http-filter-chain:bookstore/bookstore-v1": false},
			expErr:                   ErrEnvoyListenerMissing,
		},
		{
			name: "no active state listeners",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "outbound-listener",
						},
					},
				},
			},
			expectedListenerName:     "outbound-listener",
			possibleFilterChainNames: map[string]bool{"outbound-mesh-http-filter-chain:bookstore/bookstore-v1": false},
			expErr:                   ErrEnvoyActiveStateListenerMissing,
		},
		{
			name: "no filter chains in config",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "outbound-listener",
							ActiveState: &adminv3.ListenersConfigDump_DynamicListenerState{
								Listener: marshalListenerOrDie(&listenerv3.Listener{
									Name: "outbound-listener",
								}),
							},
						},
					},
				},
			},
			expectedListenerName:     "outbound-listener",
			possibleFilterChainNames: map[string]bool{"outbound-mesh-http-filter-chain:bookstore/bookstore-v1": false},
			expErr:                   ErrEnvoyFilterChainMissing,
		},
		{
			name: "no matching filter chains",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "outbound-listener",
							ActiveState: &adminv3.ListenersConfigDump_DynamicListenerState{
								Listener: marshalListenerOrDie(&listenerv3.Listener{
									Name: "outbound-listener",
									FilterChains: []*listenerv3.FilterChain{
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v1",
										},
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v2",
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedListenerName:     "outbound-listener",
			possibleFilterChainNames: map[string]bool{"outbound-mesh-tcp-filter-chain:bookstore/bookstore-v1": false},
			expErr:                   ErrEnvoyFilterChainMissing,
		},
		{
			name: "found some filter chain matches, but not for all expected filter chains",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "outbound-listener",
							ActiveState: &adminv3.ListenersConfigDump_DynamicListenerState{
								Listener: marshalListenerOrDie(&listenerv3.Listener{
									Name: "outbound-listener",
									FilterChains: []*listenerv3.FilterChain{
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v1",
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedListenerName: "outbound-listener",
			possibleFilterChainNames: map[string]bool{
				"outbound-mesh-http-filter-chain:bookstore/bookstore-v1": false,
				"outbound-mesh-http-filter-chain:bookstore/bookstore-v2": false,
			},
			expErr: ErrEnvoyFilterChainMissing,
		}, {
			name: "found all filter chain matches",
			config: &Config{
				Listeners: adminv3.ListenersConfigDump{
					DynamicListeners: []*adminv3.ListenersConfigDump_DynamicListener{
						{
							Name: "outbound-listener",
							ActiveState: &adminv3.ListenersConfigDump_DynamicListenerState{
								Listener: marshalListenerOrDie(&listenerv3.Listener{
									Name: "outbound-listener",
									FilterChains: []*listenerv3.FilterChain{
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v1",
										},
										{
											Name: "outbound-mesh-http-filter-chain:bookstore/bookstore-v2",
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedListenerName: "outbound-listener",
			possibleFilterChainNames: map[string]bool{
				"outbound-mesh-http-filter-chain:bookstore/bookstore-v1": false,
				"outbound-mesh-http-filter-chain:bookstore/bookstore-v2": false,
			},
			expErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			err := findMatchingFilterChainNames(test.config, test.expectedListenerName, test.possibleFilterChainNames)
			assert.Equal(test.expErr, err)
		})
	}
}

func marshalListenerOrDie(listener *listenerv3.Listener) *any.Any {
	a, err := ptypes.MarshalAny(listener)
	if err != nil {
		panic(err)
	}
	return a
}
