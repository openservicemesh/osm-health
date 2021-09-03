package podhelper

import (
	"testing"

	"github.com/google/uuid"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openservicemesh/osm/pkg/constants"
)

func TestHasProxyUUIDLabel(t *testing.T) {
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
					Labels: map[string]string{
						// This test requires an actual UUID
						constants.EnvoyUniqueIDLabelName: uuid.New().String()},
				},
			},
			expectedError: nil,
		},
		{
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-2",
				},
			},
			expectedError: ErrProxyUUIDLabelMissing,
		},
		{
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-3",
					Labels: map[string]string{
						// This test requires an actual UUID
						constants.EnvoyUniqueIDLabelName: "invalid-uuid"},
				},
			},
			expectedError: ErrProxyUUIDLabelMissing,
		},
	}

	for _, tc := range testCases {
		proxyUUIDLabelChecker := NewProxyUUIDLabelCheck(&tc.pod)

		assert.Equal(tc.expectedError, proxyUUIDLabelChecker.Run().GetError())
	}
}
