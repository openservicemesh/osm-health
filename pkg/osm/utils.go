package osm

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm/pkg/constants"
)

// VersionDelimiter is a delimiter used in release versions
const VersionDelimiter string = "."

// MeshInfo is the type used to represent service mesh information
type MeshInfo struct {
	Name       common.MeshName
	Namespace  common.MeshNamespace
	OSMVersion ControllerVersion
}

// GetMeshInfo returns the MeshInfo for a service mesh with its control plane in the given namespace
func GetMeshInfo(client kubernetes.Interface, osmControlPlaneNamespace string) (*MeshInfo, error) {
	osmControllerDeployment, err := GetOSMControllerDeployment(client, osmControlPlaneNamespace)
	if err != nil {
		return nil, err
	}
	osmVersion, err := FormatReleaseVersion(osmControllerDeployment.Labels[constants.OSMAppVersionLabelKey])
	if err != nil {
		return nil, err
	}

	mesh := &MeshInfo{
		Name:       common.MeshName(osmControllerDeployment.Labels[constants.OSMAppNameLabelKey]),
		Namespace:  common.MeshNamespace(osmControlPlaneNamespace),
		OSMVersion: ControllerVersion(osmVersion),
	}
	return mesh, nil
}

// GetOSMControllerDeployment returns the OSM controller deployment in a given namespace
func GetOSMControllerDeployment(client kubernetes.Interface, osmControlPlaneNamespace string) (*v1.Deployment, error) {
	deploymentsClient := client.AppsV1().Deployments(osmControlPlaneNamespace)
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
	splitVersion := strings.Split(version, VersionDelimiter)
	if len(splitVersion) < 2 {
		return "", fmt.Errorf(" is not in the expected format. The expected format is vXX.XX.XX or vXX.XX")
	}
	majorMinorVersion := splitVersion[0] + VersionDelimiter + splitVersion[1]
	return majorMinorVersion, nil
}
