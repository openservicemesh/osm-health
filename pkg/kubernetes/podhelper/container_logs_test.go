package podhelper

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHasNoBadLogs(t *testing.T) {
	type test struct {
		name                string
		pod                 corev1.Pod
		searchContainerName string
		expectedErr         error
	}

	testCases := []test{
		{
			name: "container found in container list, searching for container logs should not error",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "EnvoyContainer",
							Image: "envoyproxy/envoy-alpine:v1.18.777",
						},
						{
							Name:  "AppContainer",
							Image: "random/app:v0.0.0",
						},
					},
				},
			},
			searchContainerName: "AppContainer",
			expectedErr:         nil,
		},
		{
			name: "container found in init container list, searching for container logs should not error",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-2",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "EnvoyContainer",
							Image: "envoyproxy/envoy-alpine:v1.18.555",
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "OsmInit",
							Image: "openservicemesh/init:v0.0.0",
						},
					},
				},
			},
			searchContainerName: "OsmInit",
			expectedErr:         nil,
		},
		{
			name: "container not found, searching for the container logs should error",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-3",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "EnvoyContainer",
							Image: "envoyproxy/envoy-alpine:v1.18.555",
						},
					},
				},
			},
			searchContainerName: "RandomContainer",
			expectedErr:         ErrPodDoesNotHaveContainer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := tassert.New(t)
			client := fake.NewSimpleClientset()
			outcome := HasNoBadLogs(client, &tc.pod, tc.searchContainerName)
			assert.Equal(tc.expectedErr, outcome.GetError())
		})
	}
}
