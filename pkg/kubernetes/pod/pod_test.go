package pod

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetMatchingServices(t *testing.T) {
	type test struct {
		name                string
		podLabels           map[string]string
		namespace           string
		serviceList         []*corev1.Service
		isErrorExpected     bool
		numExpectedOutcomes int
	}

	testCases := []test{
		{
			name: "Correct matching services list returned",
			podLabels: map[string]string{
				"app":    "bookstore",
				"random": "unmatched",
			},
			namespace: "namespace-1",
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "bookstore",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "bookstore",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-3",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "not-bookstore",
						},
					},
				},
			},
			isErrorExpected:     false,
			numExpectedOutcomes: 2,
		},
		{
			name: "No matching services found, empty list returned",
			podLabels: map[string]string{
				"app":    "bookstore",
				"random": "unmatched",
			},
			namespace: "namespace-1",
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "not-bookstore",
						},
					},
				},
			},
			isErrorExpected:     false,
			numExpectedOutcomes: 0,
		},
		{
			name:      "Pod has no labels, empty list returned",
			podLabels: map[string]string{},
			namespace: "namespace-1",
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "bookstore",
						},
					},
				},
			},
			isErrorExpected:     false,
			numExpectedOutcomes: 0,
		},
		{
			name: "No services exist, empty list returned",
			podLabels: map[string]string{
				"app": "bookstore",
			},
			namespace:           "namespace-1",
			serviceList:         nil,
			isErrorExpected:     false,
			numExpectedOutcomes: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert := tassert.New(t)
			objs := make([]runtime.Object, len(testCase.serviceList))
			for i := range testCase.serviceList {
				objs[i] = testCase.serviceList[i]
			}
			client := fake.NewSimpleClientset(objs...)

			outSvc, err := GetMatchingServices(client, testCase.podLabels, testCase.namespace)
			if testCase.isErrorExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(testCase.numExpectedOutcomes, len(outSvc))
			}
		})
	}
}
