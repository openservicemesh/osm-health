package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
)

const connectivityPodToPodDesc = `
Checks connectivity between two Kubernetes pods
`

const connectivityPodToPodExample = `
Example:
	$ osm-health connectivity pod-to-pod namespace-a/pod-a namespace-b/pod-b
`

func newConnectivityPodToPodCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pod-to-pod SOURCE_POD DESTINATION_POD",
		Short:   "Checks connectivity between Kubernetes pods",
		Example: connectivityPodToPodExample,
		Long:    connectivityPodToPodDesc,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.Errorf("provide both SOURCE_POD and DESTINATION_POD")
			}

			srcPod, err := pod.FromString(args[0])
			if err != nil {
				return errors.New("invalid SOURCE_POD")
			}

			dstPod, err := pod.FromString(args[1])
			if err != nil {
				return errors.New("invalid DESTINATION_POD")
			}

			osmControlPlaneNamespace := settings.Namespace()

			connectivity.PodToPod(srcPod, dstPod, osmControlPlaneNamespace)
			return nil
		},
	}
}
