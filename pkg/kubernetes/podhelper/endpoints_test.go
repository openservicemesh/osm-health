package podhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
)

func TestEndpointsCheck(t *testing.T) {
	tests := []struct {
		name      string
		pod       *corev1.Pod
		endpoints []*corev1.Endpoints
		pass      bool
	}{
		{
			name: "ok",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "a",
					Namespace: "b",
				},
			},
			endpoints: []*corev1.Endpoints{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "a",
									},
								},
							},
						},
					},
				},
			},
			pass: true,
		},
		{
			name: "no endpoints",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "a",
					Namespace: "b",
				},
			},
			endpoints: nil,
			pass:      false,
		},
		{
			name: "endpoints in wrong namespace",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "a",
					Namespace: "b",
				},
			},
			endpoints: []*corev1.Endpoints{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "not-b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "not-b",
										Name:      "a",
									},
								},
							},
						},
					},
				},
			},
			pass: false,
		},
		{
			name: "ok when multiple Endpoints match",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "a",
					Namespace: "b",
				},
			},
			endpoints: []*corev1.Endpoints{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "a-unique-name",
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "a",
									},
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-unique-name",
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "a",
									},
								},
							},
						},
					},
				},
			},
			pass: true,
		},
		{
			name: "ok when one of several Endpoints match",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "a",
					Namespace: "b",
				},
			},
			endpoints: []*corev1.Endpoints{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "a-unique-name",
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "not-a",
									},
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "also-a-unique-name",
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "a",
									},
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-unique-name",
						Namespace: "b",
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									TargetRef: &corev1.ObjectReference{
										Namespace: "b",
										Name:      "also-not-a",
									},
								},
							},
						},
					},
				},
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objs := make([]runtime.Object, len(test.endpoints))
			for i := range test.endpoints {
				objs[i] = test.endpoints[i]
			}
			client := fake.NewSimpleClientset(objs...)
			check := NewEndpointsCheck(client, test.pod)
			out := check.Run()
			if test.pass {
				assert.Equal(t, outcomes.Pass{}, out)
			} else {
				assert.Equal(t, ErrPodNotInEndpoints, out.GetError())
			}
		})
	}
}
