package utils

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	tassert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

func TestGetMeshInfo(t *testing.T) {
	osmControlPlaneNamespace := "test-osm-system-namespace"
	osmMeshName := "test-osm-mesh-name"
	osmVersion := "v999.888.777"
	expectedOsmMajorMinorVersion := "v999.888"

	tests := []struct {
		name                  string
		deployments           []*appsv1.Deployment
		controlPlaneNamespace common.MeshNamespace
		expErr                bool
	}{
		{
			name:                  "no deployments",
			deployments:           []*appsv1.Deployment{},
			controlPlaneNamespace: common.MeshNamespace(osmControlPlaneNamespace),
			expErr:                true,
		},
		{
			name: "there are no osm-controller deployments in the desired controller namespace",
			deployments: []*appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-desired-osm-controller-deployment",
						Namespace: "not-desired-osm-namespace",
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
			},
			controlPlaneNamespace: common.MeshNamespace(osmControlPlaneNamespace),
			expErr:                true,
		},
		{
			name: "multiple deployments (controller and non-controller deployments) in the controller namespace",
			deployments: []*appsv1.Deployment{
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
			},
			controlPlaneNamespace: common.MeshNamespace(osmControlPlaneNamespace),
			expErr:                false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			objs := make([]runtime.Object, len(test.deployments))
			for i := range test.deployments {
				objs[i] = test.deployments[i]
			}
			k8s := fake.NewSimpleClientset(objs...)
			meshInfo, err := GetMeshInfo(k8s, test.controlPlaneNamespace)

			assert.Equal(test.expErr, err != nil)
			if !test.expErr {
				assert.NotNil(meshInfo)
				assert.Equal(osmMeshName, meshInfo.Name.String())
				assert.Equal(osmControlPlaneNamespace, meshInfo.Namespace.String())
				assert.Equal(expectedOsmMajorMinorVersion, meshInfo.OSMVersion.String())
			}
		})
	}
}

func TestGetOSMControllerDeployment(t *testing.T) {
	tests := []struct {
		name        string
		deployments []*appsv1.Deployment
		namespace   common.MeshNamespace
		expErr      bool
	}{
		{
			name: "there are no deployments in the controller namespace",
			deployments: []*appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "notosmcontrollerdeployment",
						Namespace: "notosmnamespace",
						Labels: map[string]string{
							constants.OSMAppInstanceLabelKey: "osm",
							constants.OSMAppVersionLabelKey:  "v0.9.1",
							"app":                            constants.OSMControllerName,
						},
					},
				},
			},
			namespace: "osmnamespace",
			expErr:    true,
		},
		{
			name: "multiple osm-controller deployments in the controller namespace",
			deployments: []*appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osmcontrollerdeployment",
						Namespace: "osmnamespace",
						Labels: map[string]string{
							constants.OSMAppInstanceLabelKey: "osm",
							constants.OSMAppVersionLabelKey:  "v0.9.1",
							"app":                            constants.OSMControllerName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osmcontrollerdeployment2",
						Namespace: "osmnamespace",
						Labels: map[string]string{
							constants.OSMAppInstanceLabelKey: "osm",
							constants.OSMAppVersionLabelKey:  "v0.9.1",
							"app":                            constants.OSMControllerName,
						},
					},
				},
			},
			namespace: "osmnamespace",
			expErr:    true,
		},
		{
			name: "multiple deployments (controller and non-controller deployments) in the controller namespace",
			deployments: []*appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osm-controller-deployment",
						Namespace: "osm-namespace",
						Labels: map[string]string{
							constants.OSMAppInstanceLabelKey: "osm",
							constants.OSMAppVersionLabelKey:  "v0.9.1",
							"app":                            constants.OSMControllerName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-osm-controller-deployment",
						Namespace: "osm-namespace",
					},
				},
			},
			namespace: "osm-namespace",
			expErr:    false,
		},
		{
			name: "there is one deployment in the controller namespace",
			deployments: []*appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osmcontrollerdeployment",
						Namespace: "osmnamespace",
						Labels: map[string]string{
							constants.OSMAppInstanceLabelKey: "osm",
							constants.OSMAppVersionLabelKey:  "v0.9.1",
							"app":                            constants.OSMControllerName,
						},
					},
				},
			},
			namespace: "osmnamespace",
			expErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			objs := make([]runtime.Object, len(test.deployments))
			for i := range test.deployments {
				objs[i] = test.deployments[i]
			}
			k8s := fake.NewSimpleClientset(objs...)
			deployment, err := GetOSMControllerDeployment(k8s, test.namespace)

			assert.Equal(test.expErr, err != nil)
			if !test.expErr {
				assert.NotNil(deployment)
				assert.Equal(test.deployments[0].Name, deployment.Name)
			}
		})
	}
}

func TestFormatReleaseVersion(t *testing.T) {
	tests := []struct {
		name                      string
		version                   string
		expectedMajorMinorVersion string
		expErr                    bool
	}{
		{
			name:                      "major, minor and patch version",
			version:                   "v0.9.0",
			expectedMajorMinorVersion: "v0.9",
			expErr:                    false,
		},
		{
			name:                      "major, minor and patch version without v prefix",
			version:                   "0.9.0",
			expectedMajorMinorVersion: "v0.9",
			expErr:                    false,
		},
		{
			name:                      "major and minor version",
			version:                   "v0.8",
			expectedMajorMinorVersion: "v0.8",
			expErr:                    false,
		},
		{
			name:                      "major and minor version without v prefix",
			version:                   "0.8",
			expectedMajorMinorVersion: "v0.8",
			expErr:                    false,
		},
		{
			name:                      "release cut version",
			version:                   "v0.8.1-rc.1",
			expectedMajorMinorVersion: "v0.8",
			expErr:                    false,
		},
		{
			name:                      "incorrectly-formatted version",
			version:                   ".1.2",
			expectedMajorMinorVersion: "",
			expErr:                    true,
		},
		{
			name:                      "major version",
			version:                   "v1",
			expectedMajorMinorVersion: "",
			expErr:                    true,
		},
		{
			name:                      "empty version",
			version:                   "",
			expectedMajorMinorVersion: "",
			expErr:                    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			majorMinorVersion, err := FormatReleaseVersion(test.version)

			assert.Equal(test.expErr, err != nil)
			assert.Equal(test.expectedMajorMinorVersion, majorMinorVersion)
		})
	}
}

func TestGetMonitoredNamespaces(t *testing.T) {
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
		name                            string
		namespaces                      []*corev1.Namespace
		expectedMonitoredNamespaceNames mapset.Set
	}{
		{
			name:                            "no namespaces in cluster",
			namespaces:                      nil,
			expectedMonitoredNamespaceNames: mapset.NewSet(),
		},
		{
			name: "no monitored namespaces in cluster",
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-monitored-namespace-1",
						Namespace: "not-monitored-namespace-1",
					},
				},
			},
			expectedMonitoredNamespaceNames: mapset.NewSet(),
		},
		{
			name: "multiple namespaces (monitored and unmonitored) in cluster",
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
			expectedMonitoredNamespaceNames: mapset.NewSet("monitored-namespace-1", "monitored-namespace-2"),
		},
		{
			name: "multiple namespaces (monitored and ignored) in cluster",
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace",
						Namespace: "monitored-namespace",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ignored-namespace",
						Namespace: "ignored-namespace",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
							// As of v0.10, the mutating webhook ONLY checks for the existence of the "constants.IgnoreLabel" label key.
							// As long as the key is present in the namespace labels, the namespace will be ignored. It does not check the value.
							constants.IgnoreLabel: "true",
						},
					},
				},
			},
			expectedMonitoredNamespaceNames: mapset.NewSet("monitored-namespace"),
		},
		{
			name: "multiple namespaces (monitored namespace and osm mesh namespace) in cluster",
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace",
						Namespace: "monitored-namespace",
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      osmControlPlaneNamespace,
						Namespace: osmControlPlaneNamespace,
						Labels: map[string]string{
							constants.OSMKubeResourceMonitorAnnotation: osmMeshName,
						},
					},
				},
			},
			expectedMonitoredNamespaceNames: mapset.NewSet("monitored-namespace"),
		},
		{
			name: "multiple namespaces (monitored namespace and namespace with control-plane label key) in cluster",
			namespaces: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monitored-namespace",
						Namespace: "monitored-namespace",
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
							"control-plane": "key-exists",
						},
					},
				},
			},
			expectedMonitoredNamespaceNames: mapset.NewSet("monitored-namespace"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)

			var objs []runtime.Object
			for _, ns := range test.namespaces {
				objs = append(objs, ns)
			}
			for _, deployment := range osmControlPlaneNamespaceDeployments {
				objs = append(objs, deployment)
			}
			client := fake.NewSimpleClientset(objs...)

			monitoredNamespaces, err := GetMonitoredNamespaces(client, common.MeshNamespace(osmControlPlaneNamespace))
			assert.Nil(err)
			assert.Equal(test.expectedMonitoredNamespaceNames.Cardinality(), len(monitoredNamespaces.Items))
			for _, ns := range monitoredNamespaces.Items {
				assert.True(test.expectedMonitoredNamespaceNames.Contains(ns.Name))
			}
		})
	}
}
