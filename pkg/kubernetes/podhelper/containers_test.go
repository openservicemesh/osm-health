package podhelper

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHasExpectedNumContainers(t *testing.T) {
	assert := tassert.New(t)

	type test struct {
		pod           corev1.Pod
		expectedError error
	}

	testCases := []test{
		{
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
			expectedError: nil,
		},
		{
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
			expectedError: ErrExpectedMinNumContainers,
		},
		{
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
			expectedError: ErrExpectedMinNumContainers,
		},
	}

	for _, tc := range testCases {
		numContainersChecker := HasMinExpectedContainers(&tc.pod, 2)

		assert.Equal(tc.expectedError, numContainersChecker.Run())
	}
}
