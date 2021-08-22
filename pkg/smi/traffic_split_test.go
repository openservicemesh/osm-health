package smi

import (
	"testing"

	split "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/split/v1alpha2"
	fakeSmiSplitClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned/fake"
	tassert "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestIsInTrafficSplit(t *testing.T) {
	type test struct {
		name                     string
		pod                      corev1.Pod
		serviceList              []*corev1.Service
		trafficSplitList         []*split.TrafficSplit
		isErrorExpected          bool
		isDiagnosticInfoExpected bool
	}

	testCases := []test{
		{
			name: "Matching service and split found, no error",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "namespace-1",
					Labels: map[string]string{
						"app": "bookstore-v1",
					},
				},
			},
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "bookstore-v1",
						},
					},
				},
			},
			trafficSplitList: []*split.TrafficSplit{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "traffic-split-1",
						Namespace: "namespace-1",
					},
					Spec: split.TrafficSplitSpec{
						Service: "svc",
						Backends: []split.TrafficSplitBackend{
							{
								Service: "service-1",
							},
						},
					},
				},
			},
			isErrorExpected:          false,
			isDiagnosticInfoExpected: false,
		},
		{
			name: "No split found referring to the pod's service, error not returned, diagnostic info expected",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "namespace-1",
					Labels: map[string]string{
						"app": "bookstore-v1",
					},
				},
			},
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "bookstore-v1",
						},
					},
				},
			},
			trafficSplitList: []*split.TrafficSplit{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "traffic-split-1",
						Namespace: "namespace-1",
					},
					Spec: split.TrafficSplitSpec{
						Service: "svc",
						Backends: []split.TrafficSplitBackend{
							{
								Service: "not-service-1",
							},
						},
					},
				},
			},
			isErrorExpected:          false,
			isDiagnosticInfoExpected: true,
		},
		{
			name: "No matching svc found for pod, error not returned, diagnostic info expected",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-1",
					Namespace: "namespace-1",
					Labels: map[string]string{
						"app": "bookstore-v1",
					},
				},
			},
			serviceList: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "namespace-1",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{
							"app": "not-bookstore-v1",
						},
					},
				},
			},
			trafficSplitList: []*split.TrafficSplit{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "traffic-split-1",
						Namespace: "namespace-1",
					},
					Spec: split.TrafficSplitSpec{
						Service: "svc",
						Backends: []split.TrafficSplitBackend{
							{
								Service: "service-1",
							},
						},
					},
				},
			},
			isErrorExpected:          false,
			isDiagnosticInfoExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert := tassert.New(t)
			svcObjs := make([]runtime.Object, len(testCase.serviceList))
			for i := range testCase.serviceList {
				svcObjs[i] = testCase.serviceList[i]
			}
			client := fake.NewSimpleClientset(svcObjs...)

			splitObjs := make([]runtime.Object, len(testCase.trafficSplitList))
			for i := range testCase.trafficSplitList {
				splitObjs[i] = testCase.trafficSplitList[i]
			}
			smiSplitClient := fakeSmiSplitClient.NewSimpleClientset(splitObjs...)

			trafficSplitChecker := IsInTrafficSplit(client, &testCase.pod, smiSplitClient)
			if testCase.isErrorExpected {
				assert.Error(trafficSplitChecker.Run().GetError())
			} else {
				assert.NoError(trafficSplitChecker.Run().GetError())
			}
			if testCase.isDiagnosticInfoExpected {
				assert.NotEmpty(trafficSplitChecker.Run().GetLongDiagnostics())
			}
		})
	}
}
