package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/cli"
	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
)

const connectivityPodToPodDesc = `
Checks connectivity between two Kubernetes pods
`

const connectivityPodToPodExample = `$ osm-health connectivity pod-to-pod source-namespace/source-pod destination-namespace/destination-pod`

func newConnectivityPodToPodCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pod-to-pod source-namespace/source-pod destination-namespace/destination-pod",
		Short:   "Checks connectivity between Kubernetes pods",
		Example: connectivityPodToPodExample,
		Long:    connectivityPodToPodDesc,
		Args:    cli.ExactArgsWithError(2, errors.New("requires 2 arguments: source-namespace/source-pod destination-namespace/destination-pod")),
		RunE: func(_ *cobra.Command, args []string) error {
			srcPod, err := pod.FromString(args[0])
			if err != nil {
				return errors.New("invalid source-namespace/source-pod")
			}

			dstPod, err := pod.FromString(args[1])
			if err != nil {
				return errors.New("invalid destination-namespace/destination-pod")
			}

			osmControlPlaneNamespace := settings.Namespace()

			connectivity.PodToPod(srcPod, dstPod, osmControlPlaneNamespace)
			return nil
		},
	}
}
