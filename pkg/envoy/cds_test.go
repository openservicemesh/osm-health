package envoy

import (
	"testing"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEnvoyClusterChecker(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		dstPod *corev1.Pod
		svcs   []*corev1.Service
		pass   bool
	}{
		{
			name: "pod matches one service with cluster in config",
			config: &Config{
				Clusters: adminv3.ClustersConfigDump{
					DynamicActiveClusters: []*adminv3.ClustersConfigDump_DynamicCluster{
						{
							Cluster: marshalClusterOrDie(&clusterv3.Cluster{
								Name: "mynamespace/myservice",
							}),
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
			pass: true,
		},
		{
			name: "pod matches one service without cluster in config",
			config: &Config{
				Clusters: adminv3.ClustersConfigDump{
					DynamicActiveClusters: []*adminv3.ClustersConfigDump_DynamicCluster{
						{
							Cluster: marshalClusterOrDie(&clusterv3.Cluster{
								Name: "mynamespace/not-myservice",
							}),
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
			pass: false,
		},
		{
			name: "pod matches two services with one cluster in config",
			config: &Config{
				Clusters: adminv3.ClustersConfigDump{
					DynamicActiveClusters: []*adminv3.ClustersConfigDump_DynamicCluster{
						{
							Cluster: marshalClusterOrDie(&clusterv3.Cluster{
								Name: "mynamespace/myservice",
							}),
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
			pass: true,
		},
		{
			name:   "pod matches no services with no clusters",
			config: &Config{},
			dstPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "mynamespace",
				},
			},
			svcs: nil,
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
			objs := make([]runtime.Object, len(test.svcs))
			for i := range test.svcs {
				objs[i] = test.svcs[i]
			}
			k8s := fake.NewSimpleClientset(objs...)
			clusterChecker := HasCluster(k8s, configGetter, test.dstPod)
			err := clusterChecker.Run()
			if test.pass {
				assert.NoError(err)
			} else {
				assert.Error(err)
			}
		})
	}
}

func marshalClusterOrDie(cluster *clusterv3.Cluster) *any.Any {
	a, err := ptypes.MarshalAny(cluster)
	if err != nil {
		panic(err)
	}
	return a
}
