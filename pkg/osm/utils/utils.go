package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/osm/version"
	"github.com/openservicemesh/osm/pkg/constants"
)

const versionDelimiter string = "."

// MeshInfo is the type used to represent service mesh information
type MeshInfo struct {
	Name       common.MeshName
	Namespace  common.MeshNamespace
	OSMVersion version.ControllerVersion
}

// GetMeshInfo returns the MeshInfo for a service mesh with its control plane in the given namespace
func GetMeshInfo(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace) (*MeshInfo, error) {
	osmControllerDeployment, err := GetOSMControllerDeployment(client, osmControlPlaneNamespace)
	if err != nil {
		return nil, err
	}
	osmVersion, err := FormatReleaseVersion(osmControllerDeployment.Labels[constants.OSMAppVersionLabelKey])
	if err != nil {
		return nil, err
	}

	mesh := &MeshInfo{
		Name:       common.MeshName(osmControllerDeployment.Labels[constants.OSMAppInstanceLabelKey]),
		Namespace:  osmControlPlaneNamespace,
		OSMVersion: version.ControllerVersion(osmVersion),
	}
	return mesh, nil
}

// GetOSMControllerDeployment returns the OSM controller deployment in a given namespace
func GetOSMControllerDeployment(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace) (*v1.Deployment, error) {
	deploymentsClient := client.AppsV1().Deployments(osmControlPlaneNamespace.String())
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": constants.OSMControllerName}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	deployments, err := deploymentsClient.List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	if len(deployments.Items) == 0 {
		return nil, fmt.Errorf("%s deployment not found in %s namespace", constants.OSMControllerName, osmControlPlaneNamespace)
	} else if len(deployments.Items) > 1 {
		return nil, fmt.Errorf("found more than one %s deployments in %s namespace", constants.OSMControllerName, osmControlPlaneNamespace)
	}
	return &deployments.Items[0], nil
}

// FormatReleaseVersion returns the major and minor version of the release
func FormatReleaseVersion(version string) (string, error) {
	splitVersion := strings.Split(version, versionDelimiter)
	if len(splitVersion) < 2 {
		return "", fmt.Errorf("%s is not in the expected format", version)
	}
	majorMinorVersion := splitVersion[0] + versionDelimiter + splitVersion[1]
	return majorMinorVersion, nil
}

// GetMonitoredNamespaces returns a list of namespaces monitored by the osm control plane.
func GetMonitoredNamespaces(client kubernetes.Interface, osmControlPlaneNamespace common.MeshNamespace) (*corev1.NamespaceList, error) {
	meshInfo, err := GetMeshInfo(client, osmControlPlaneNamespace)
	if err != nil {
		return &corev1.NamespaceList{}, err
	}
	meshName := meshInfo.Name

	// Criteria that is used to determine if a namespace is monitored by osm (from osm's mutating webhook):
	// 		kubectl get MutatingWebhookConfiguration osm-webhook-osm -o json | jq '.webhooks[0].namespaceSelector'
	/*
		{
		  "matchExpressions": [
		    {
		      "key": "openservicemesh.io/ignore",
		      "operator": "DoesNotExist"
		    },
		    {
		      "key": "name",
		      "operator": "NotIn",
		      "values": [
		        "osm-system"
		      ]
		    }
		  ],
		  "matchLabels": {
		    "openservicemesh.io/monitored-by": "osm" // <<--- value should be the mesh name
		  }
		}
	*/

	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			constants.OSMKubeResourceMonitorAnnotation: meshName.String(),
		},
	}

	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}

	nsList, err := client.CoreV1().Namespaces().List(context.TODO(), listOptions)
	if err != nil {
		return &corev1.NamespaceList{}, errors.Errorf("unable to list namespaces in the cluster: %s", err.Error())
	}

	var monitoredNamespaces []corev1.Namespace
	for _, ns := range nsList.Items {
		if _, exists := ns.GetLabels()[constants.IgnoreLabel]; exists {
			// As of v0.10, the mutating webhook ONLY checks for the existence of the "constants.IgnoreLabel" label key.
			// As long as the key is present in the namespace labels, the namespace will be ignored. It does not check the value.
			continue
		}
		if ns.Name == osmControlPlaneNamespace.String() {
			continue
		}
		if _, exists := ns.GetLabels()["control-plane"]; exists {
			continue
		}
		monitoredNamespaces = append(monitoredNamespaces, ns)
	}

	return &corev1.NamespaceList{Items: monitoredNamespaces}, nil
}
