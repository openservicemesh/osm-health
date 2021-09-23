package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	tassert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

func TestCheckControllerProxyConnectionMetrics(t *testing.T) {
	var osmProxyConnectCountMetricID = "osm_proxy_connect_count"
	osmControlPlaneNamespace := "test-osm-system-namespace"
	osmMeshName := "test-osm-mesh-name"
	osmVersion := "v999.888.777"

	osmControlPlaneNamespaceDeployments := []*appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "osm-controller-deployment",
				Namespace: osmControlPlaneNamespace,
				Labels: map[string]string{
					constants.OSMAppInstanceLabelKey: osmMeshName,
					constants.OSMAppVersionLabelKey:  osmVersion,
					"app":                            constants.OSMControllerName,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "not-osm-controller-deployment",
				Namespace: osmControlPlaneNamespace,
			},
		},
	}

	tests := []struct {
		name                      string
		statusCode                int
		httpServerMetricsRespBody string
		namespaces                []*corev1.Namespace
		pods                      []*corev1.Pod
		expectedError             error
	}{
		{
			name:                      "http metrics server returns service unavailable http status code",
			statusCode:                http.StatusServiceUnavailable,
			httpServerMetricsRespBody: "",
			namespaces:                nil,
			pods:                      nil,
			expectedError:             errors.New("osm-controller metrics check failed: url returned HTTP status code: 503"),
		},
		{
			name:                      "error: http server returns 3 pods, but 2 pods in 2 monitored namespaces, 1 pod in unmonitored namespace.",
			statusCode:                http.StatusOK,
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 3\nrandom", osmProxyConnectCountMetricID),
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-1",
						Namespace: "monitored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-2",
						Namespace: "monitored-namespace-2",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-monitored-namespace-1",
						Namespace: "not-monitored-namespace-1",
					},
				},
			},
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-1",
						Namespace: "monitored-namespace-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-2",
						Namespace: "monitored-namespace-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-monitored-pod-1",
						Namespace: "not-monitored-namespace-1",
					},
				},
			},
			expectedError: errors.Errorf("osm-controller metrics check failed: incorrect %s metric: expected 2 but http server metrics returned 3", osmProxyConnectCountMetricID),
		},
		{
			name:                      "no error: correct metric returned, http server returns 2, 2 pods in 2 monitored namespaces. 1 pod in unmonitored namespace.",
			statusCode:                http.StatusOK,
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 2\nrandom", osmProxyConnectCountMetricID),
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-1",
						Namespace: "monitored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-2",
						Namespace: "monitored-namespace-2",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-monitored-namespace-1",
						Namespace: "not-monitored-namespace-1",
					},
				},
			},
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-1",
						Namespace: "monitored-namespace-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-2",
						Namespace: "monitored-namespace-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-monitored-pod-1",
						Namespace: "not-monitored-namespace-1",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:                      "no error: correct metric returned, http server returns 2, 2 pods in 2 monitored namespaces. 1 pod in ignored namespace.",
			statusCode:                http.StatusOK,
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 2\nrandom", osmProxyConnectCountMetricID),
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-1",
						Namespace: "monitored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-2",
						Namespace: "monitored-namespace-2",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ignored-namespace-1",
						Namespace: "ignored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
							// As of v0.10, the mutating webhook ONLY checks for the existence of the "constants.IgnoreLabel" label key.
							// As long as the key is present in the namespace labels, the namespace will be ignored. It does not check the value.
							constants.IgnoreLabel: "true",
						},
					},
				},
			},
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-1",
						Namespace: "monitored-namespace-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-2",
						Namespace: "monitored-namespace-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ignored-pod-1",
						Namespace: "ignored-namespace-1",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:                      "no error: correct metric returned, http server returns 2, 2 pods in 2 monitored namespaces. 1 pod in osm control plane namespace.",
			statusCode:                http.StatusOK,
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 2\nrandom", osmProxyConnectCountMetricID),
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-1",
						Namespace: "monitored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-2",
						Namespace: "monitored-namespace-2",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
			},
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-1",
						Namespace: "monitored-namespace-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-2",
						Namespace: "monitored-namespace-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osm-namespace-pod",
						Namespace: osmControlPlaneNamespace,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:                      "no error: correct metric returned, http server returns 2, 2 pods in 2 monitored namespaces. 1 pod in namespace with 'control-plane' label.",
			statusCode:                http.StatusOK,
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 2\nrandom", osmProxyConnectCountMetricID),
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-1",
						Namespace: "monitored-namespace-1",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace-2",
						Namespace: "monitored-namespace-2",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "control-plane-namespace",
						Namespace: "control-plane-namespace",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
							// As of v0.10, the mutating webhook ONLY checks for the existence of the "control-plane" label key.
							// As long as the key is present in the namespace labels, the namespace will be ignored. It does not check the value.
							"control-plane": "true",
						},
					},
				},
			},
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-1",
						Namespace: "monitored-namespace-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-pod-2",
						Namespace: "monitored-namespace-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "control-plane-pod",
						Namespace: "control-plane-namespace",
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				_, _ = w.Write([]byte(test.httpServerMetricsRespBody))
			}))
			defer ts.Close()

			var objs []runtime.Object
			for _, ns := range test.namespaces {
				objs = append(objs, ns)
			}
			for _, deployment := range osmControlPlaneNamespaceDeployments {
				objs = append(objs, deployment)
			}
			client := fake.NewSimpleClientset(objs...)
			for _, pod := range test.pods {
				_, err := client.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
				if err != nil {
					log.Panic().Msgf("unable to create pod %s in namespace %s in test %s", pod.Name, pod.Namespace, test.name)
				}
			}

			err := checkControllerProxyConnectionMetrics(client, ts.URL, common.MeshNamespace(osmControlPlaneNamespace))

			assert.Equal(test.expectedError != nil, err != nil)
			if test.expectedError != nil {
				assert.Equal(test.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestCheckProxyConnectCount(t *testing.T) {
	var osmProxyConnectCountMetricID = "osm_proxy_connect_count"
	tests := []struct {
		name                      string
		httpServerMetricsRespBody string
		expectedError             error
		expectedProxyConnectCount int
	}{
		{
			name:                      fmt.Sprintf("valid %s metric in http server metrics response body", osmProxyConnectCountMetricID),
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 55\nrandom", osmProxyConnectCountMetricID),
			expectedError:             nil,
			expectedProxyConnectCount: 55,
		},
		{
			name:                      fmt.Sprintf("incorrect %s metric value in http server metrics response body", osmProxyConnectCountMetricID),
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s 19\nrandom", osmProxyConnectCountMetricID),
			expectedError: errors.Errorf("incorrect %s metric: expected %d but http server metrics returned %d",
				osmProxyConnectCountMetricID, 55, 19),
			expectedProxyConnectCount: 55,
		},
		{
			name:                      "empty http server metrics response body",
			httpServerMetricsRespBody: "",
			expectedError:             errors.Errorf("missing or invalid %s metric in HTTP server metrics response", osmProxyConnectCountMetricID),
			expectedProxyConnectCount: 0,
		},
		{
			name:                      fmt.Sprintf("missing %s metric in http server metrics response body", osmProxyConnectCountMetricID),
			httpServerMetricsRespBody: "random",
			expectedError:             errors.Errorf("missing or invalid %s metric in HTTP server metrics response", osmProxyConnectCountMetricID),
			expectedProxyConnectCount: 0,
		},
		{
			name:                      fmt.Sprintf("invalid %s metric in http server metrics response body", osmProxyConnectCountMetricID),
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s ABC\nrandom", osmProxyConnectCountMetricID),
			expectedError:             errors.Errorf("missing or invalid %s metric in HTTP server metrics response", osmProxyConnectCountMetricID),
			expectedProxyConnectCount: 0,
		},
		{
			name:                      fmt.Sprintf("invalid %s metric in http server metrics response body", osmProxyConnectCountMetricID),
			httpServerMetricsRespBody: fmt.Sprintf("random\n%s ABC\nrandom", osmProxyConnectCountMetricID),
			expectedError:             errors.Errorf("missing or invalid %s metric in HTTP server metrics response", osmProxyConnectCountMetricID),
			expectedProxyConnectCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			err := checkProxyConnectCount(test.expectedProxyConnectCount, test.httpServerMetricsRespBody)
			assert.Equal(test.expectedError != nil, err != nil)
			if test.expectedError != nil {
				assert.Equal(test.expectedError.Error(), err.Error())
			}
		})
	}
}
