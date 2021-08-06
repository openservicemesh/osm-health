package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

const connectivityPodToPodDesc = `
Checks connectivity between Kubernetes pods
	(add more descriptive description)
`

const connectivityPodToPodExample = `
Example:
	(add example)
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

			fromPod, err := kuberneteshelper.PodFromString(args[0])
			if err != nil {
				return errors.New("invalid SOURCE_POD")
			}

			toPod, err := kuberneteshelper.PodFromString(args[1])
			if err != nil {
				return errors.New("invalid DESTINATION_POD")
			}

			osmControlPlaneNamespace := settings.Namespace()

			connectivity.PodToPod(fromPod, toPod, osmControlPlaneNamespace)
			return nil
		},
	}
}
