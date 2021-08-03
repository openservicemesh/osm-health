package connectivity

import (
	"net/url"

	v1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm-health/pkg/common"
)

// PodToURL tests the connectivity between a source pod and destination url.
func PodToURL(fromPod *v1.Pod, destinationURL *url.URL) common.Result {
	log.Info().Msgf("Testing connectivity from %s/%s to %s", fromPod.Namespace, fromPod.Name, destinationURL)

	outcomes := common.Run()

	common.Print(outcomes...)

	return common.Result{
		SMIPolicy: common.SMIPolicy{
			HasPolicy:                  false,
			ValidPolicy:                false,
			SourceToDestinationAllowed: false,
		},
		Successful: false,
	}
}
