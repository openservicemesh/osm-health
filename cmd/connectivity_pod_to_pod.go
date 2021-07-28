package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"

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
	cmd := &cobra.Command{
		Use:   "pod-to-pod SOURCE_POD DESTINATION_POD",
		Short: "Checks connectivity between Kubernetes pods",
		Long:  connectivityPodToPodDesc,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.Errorf("provide both SOURCE_POD and DESTINATION_POD")
			}

			config, err := settings.RESTClientGetter().ToRESTConfig()
			if err != nil {
				return errors.Errorf("Error fetching kubeconfig: %s", err)
			}

			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				return errors.Errorf("Could not access Kubernetes cluster, check kubeconfig: %s", err)
			}
			clientSet := clientset

			fromPod, err := kuberneteshelper.PodFromString(args[0])
			if err != nil {
				return errors.New("invalid SOURCE_POD")
			}

			toPod, err := kuberneteshelper.PodFromString(args[1])
			if err != nil {
				return errors.New("invalid DESTINATION_POD")
			}

			connectivity.PodToPod(clientSet, fromPod, toPod)
			return nil
		},
		Example: connectivityPodToPodExample,
	}
	return cmd
}
