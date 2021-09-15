package envoy

import (
	"fmt"
	"testing"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm-health/pkg/runner"
)

func TestEnvoySecretCheck(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func(kubernetes.Interface, ConfigGetter, *corev1.Pod) runner.Runnable
		config    *Config
		dstPod    *corev1.Pod
		svcs      []*corev1.Service
		svcAccts  []*corev1.ServiceAccount
		pass      bool
	}{
		{
			name:      "pod matches one service with outbound mTLS root certificate secret in config",
			checkFunc: HasOutboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSOutbound, "mynamespace/myservice"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
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
			},
			svcAccts: nil,
			pass:     true,
		},
		{
			name:      "pod matches one service without outbound mTLS root certificate secret in config",
			checkFunc: HasOutboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSOutbound, "mynamespace/not-myservice"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
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
			},
			svcAccts: nil,
			pass:     false,
		},
		{
			name:      "pod matches two services with one outbound mTLS root certificate secret in config",
			checkFunc: HasOutboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSOutbound, "mynamespace/myservice"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
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
			svcAccts: nil,
			pass:     true,
		},
		{
			name:      "pod matches no services with no outbound mTLS root certificate secrets",
			checkFunc: HasOutboundRootCertificate,
			config:    &Config{},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
				},
			},
			svcs:     nil,
			svcAccts: nil,
			pass:     false,
		},
		{
			name:      "pod has serviceaccount and inbound mTLS root certificate secret in config",
			checkFunc: HasInboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSInbound, "mynamespace/myserviceaccount"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "myserviceaccount",
				},
			},
			svcs: nil,
			svcAccts: []*corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myserviceaccount",
						Namespace: "mynamespace",
					},
				},
			},
			pass: true,
		},
		{
			name:      "pod has serviceaccount and no inbound mTLS root certificate secret in config",
			checkFunc: HasInboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSInbound, "mynamespace/not-myserviceaccount"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "myserviceaccount",
				},
			},
			svcs: nil,
			svcAccts: []*corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myserviceaccount",
						Namespace: "mynamespace",
					},
				},
			},
			pass: false,
		},
		{
			name:      "pod has serviceaccount and no inbound mTLS root certificate secrets in config",
			checkFunc: HasInboundRootCertificate,
			config:    &Config{},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "myserviceaccount",
				},
			},
			svcs: nil,
			svcAccts: []*corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myserviceaccount",
						Namespace: "mynamespace",
					},
				},
			},
			pass: false,
		},
		{
			name:      "pod has no serviceaccount and no inbound mTLS root certificate secret in config",
			checkFunc: HasInboundRootCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", RootCertTypeForMTLSInbound, "mynamespace/not-myserviceaccount"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
			},
			svcs:     nil,
			svcAccts: nil,
			pass:     false,
		},
		{
			name:      "pod has serviceaccount and service certificate secret in config",
			checkFunc: HasServiceCertificate,
			config: &Config{
				SecretsConfigDump: adminv3.SecretsConfigDump{
					DynamicActiveSecrets: []*adminv3.SecretsConfigDump_DynamicSecret{
						{
							Name: fmt.Sprintf("%s:%s", ServiceCertType, "mynamespace/myserviceaccount"),
						},
					},
				},
			},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "myserviceaccount",
				},
			},
			svcs: nil,
			svcAccts: []*corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myserviceaccount",
						Namespace: "mynamespace",
					},
				},
			},
			pass: true,
		},
		{
			name:      "pod has no serviceaccount and no service certificate secret in config",
			checkFunc: HasServiceCertificate,
			config:    &Config{},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
					Labels: map[string]string{
						"mykey": "myval",
					},
				},
			},
			svcs:     nil,
			svcAccts: nil,
			pass:     false,
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
			objs := make([]runtime.Object, len(test.svcs))
			for i := range test.svcs {
				objs[i] = test.svcs[i]
			}
			k8s := fake.NewSimpleClientset(objs...)
			envoyCertificateChecker := test.checkFunc(k8s, configGetter, test.dstPod)
			outcome := envoyCertificateChecker.Run()
			if test.pass {
				assert.NoError(outcome.GetError())
			} else {
				assert.Error(outcome.GetError())
			}
		})
	}
}
