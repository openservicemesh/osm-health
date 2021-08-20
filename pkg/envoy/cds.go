package envoy

import (
	"fmt"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/common"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
	"github.com/openservicemesh/osm/pkg/utils"
)

// HasClusterCheck implements common.Runnable
type HasClusterCheck struct {
	ConfigGetter

	dstPod *corev1.Pod
	k8s    kubernetes.Interface
}

// Run implements common.Runnable
func (c HasClusterCheck) Run() error {
	if c.ConfigGetter == nil {
		log.Error().Msg("Incorrectly initialized ConfigGetter")
		return ErrIncorrectlyInitializedConfigGetter
	}
	envoyConfig, err := c.ConfigGetter.GetConfig()
	if err != nil {
		return err
	}

	if envoyConfig == nil {
		return ErrEnvoyConfigEmpty
	}

	// The destination Pod might back multiple services, so check that at least
	// one of those services is listed as a cluster in the source Envoy config.
	possibleClusterNames := map[string]struct{}{}
	svcs, err := kuberneteshelper.GetMatchingServices(c.k8s, c.dstPod.Labels, c.dstPod.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to map Pod %s/%s to Kubernetes Services", c.dstPod.Namespace, c.dstPod.Name)
	}
	for _, svc := range svcs {
		possibleClusterNames[utils.K8sSvcToMeshSvc(svc).String()] = struct{}{}
	}
	if len(possibleClusterNames) == 0 {
		// This pod isn't backing any services, so we wouldn't expect a cluster
		// to be listed in the Envoy config.
		return nil
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

		if _, exists := possibleClusterNames[cluster.Name]; exists {
			found = true
			break
		}
	}

	if !found {
		var expectedClusterNames []string
		for name := range possibleClusterNames {
			expectedClusterNames = append(expectedClusterNames, name)
		}
		return fmt.Errorf("Expected a cluster named one of %v, but only found %v", expectedClusterNames, foundClusterNames)
	}
	return nil
}

// Suggestion implements common.Runnable
func (c HasClusterCheck) Suggestion() string {
	panic("implement me")
}

// FixIt implements common.Runnable
func (c HasClusterCheck) FixIt() error {
	panic("implement me")
}

// Info implements common.Runnable
func (c HasClusterCheck) Info() string {
	return fmt.Sprintf("Checking whether %s is configured with an envoy cluster referring to Pod %s/%s", c.ConfigGetter.GetObjectName(), c.dstPod.Namespace, c.dstPod.Name)
}

// HasCluster creates a new common.Runnable, which checks whether the given Pod
// has an Envoy with properly configured cluster.
func HasCluster(client kubernetes.Interface, configGetter ConfigGetter, dstPod *corev1.Pod) common.Runnable {
	return HasClusterCheck{
		ConfigGetter: configGetter,
		dstPod:       dstPod,
		k8s:          client,
	}
}
