package osm

import (
	"testing"

	tassert "github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

func TestGetOSMControllerDeployment(t *testing.T) {
	tests := []struct {
		name        string
		deployments []*v1.Deployment
		namespace   common.MeshNamespace
		expErr      bool
	}{
		{
			name: "there are no deployments in the controller namespace",
			deployments: []*v1.Deployment{
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
			name: "multiple deployments in the controller namespace",
			deployments: []*v1.Deployment{
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
			name: "there is one deployment in the controller namespace",
			deployments: []*v1.Deployment{
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
			name:                      "major and minor version",
			version:                   "v0.8",
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
