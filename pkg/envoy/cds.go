package envoy

import (
	"fmt"
	"strings"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common/outcomes"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
	"github.com/openservicemesh/osm-health/pkg/runner"
	"github.com/openservicemesh/osm-health/pkg/utils"
)

// Verify interface compliance
var _ runner.Runnable = (*ClusterCheck)(nil)

// ClusterCheck implements common.Runnable
type ClusterCheck struct {
	ConfigGetter

	dstPod *corev1.Pod
	k8s    kubernetes.Interface
}

// Run implements common.Runnable
func (c ClusterCheck) Run() outcomes.Outcome {
	if c.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return outcomes.Fail{Error: ErrIncorrectlyInitializedConfigGetter}
	}
	envoyConfig, err := c.ConfigGetter.GetConfig()
	if err != nil {
		return outcomes.Fail{Error: err}
	}

	if envoyConfig == nil {
		return outcomes.Fail{Error: ErrEnvoyConfigEmpty}
	}

	// The destination Pod might back multiple services, so check that at least
	// one of those services is listed as a cluster in the source Envoy config.
	possibleClusterNames := map[string]struct{}{}
	svcs, err := pod.GetMatchingServices(c.k8s, c.dstPod.Labels, c.dstPod.Namespace)
	if err != nil {
		return outcomes.Fail{Error: errors.Wrapf(err, "failed to map Pod %s/%s to Kubernetes Services", c.dstPod.Namespace, c.dstPod.Name)}
	}
	for _, svc := range svcs {
		possibleClusterNames[utils.K8sSvcToMeshSvc(svc).String()] = struct{}{}
	}
	if len(possibleClusterNames) == 0 {
		// This pod isn't backing any services, so we wouldn't expect a cluster
		// to be listed in the Envoy config.
		return outcomes.Info{Diagnostics: "pod is not backing any services - no clusters listed in the Envoy config"}
	}

	found := false
	var foundClusterNames []string
	for _, dynCluster := range envoyConfig.Clusters.DynamicActiveClusters {
		var cluster clusterv3.Cluster
		err := dynCluster.Cluster.UnmarshalTo(&cluster)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unmarshal cluster %s", dynCluster.String())
			continue
		}
		foundClusterNames = append(foundClusterNames, cluster.Name)

		// Beginning in OSM version v0.10 onwards, the Envoy cluster name is appended with the port number of the service.
		// cluster.Name pre v0.9 would look like "bookstore/bookstore-v1"
		// cluster.Name v0.10 onwards would look like "bookstore/bookstore-v1|14001"
		splitEnvoyClusterName := strings.Split(cluster.Name, "|")
		if len(splitEnvoyClusterName) == 0 || len(splitEnvoyClusterName[0]) == 0 {
			continue
		}
		if _, exists := possibleClusterNames[splitEnvoyClusterName[0]]; exists {
			found = true
			break
		}
	}

	if !found {
		var expectedClusterNames []string
		for name := range possibleClusterNames {
			expectedClusterNames = append(expectedClusterNames, name)
		}
		return outcomes.Fail{Error: fmt.Errorf("Expected a cluster named one of %v, but only found %v", expectedClusterNames, foundClusterNames)}
	}
	return outcomes.Pass{}
}

// Suggestion implements common.Runnable
func (c ClusterCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (c ClusterCheck) FixIt() error {
	panic("implement me")
}

// Description implements common.Runnable
func (c ClusterCheck) Description() string {
	return fmt.Sprintf("Checking whether %s is configured with an envoy cluster referring to Pod %s/%s", c.ConfigGetter.GetObjectName(), c.dstPod.Namespace, c.dstPod.Name)
}

// NewClusterCheck creates a ClusterCheck which checks whether the given Pod has an Envoy with properly configured cluster.
func NewClusterCheck(client kubernetes.Interface, configGetter ConfigGetter, dstPod *corev1.Pod) ClusterCheck {
	return ClusterCheck{
		ConfigGetter: configGetter,
		dstPod:       dstPod,
		k8s:          client,
	}
}
