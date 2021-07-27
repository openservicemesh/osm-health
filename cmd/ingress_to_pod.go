package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm-health/pkg/ingress"
	"github.com/openservicemesh/osm-health/pkg/kuberneteshelper"
)

func newIngressToPodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to-pod DESTINATION_POD",
		Short: "Checks ingress to a given Kubernetes pod",
		Long:  `Checks ingress to a given Kubernetes pod`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.Errorf("missing DESTINATION_POD parameter")
			}
			log.Info().Msgf("Checking Ingress to Pod %s", args[0])
			config, err := settings.RESTClientGetter().ToRESTConfig()
			if err != nil {
				return errors.Errorf("Error fetching kubeconfig: %s", err)
			}

			clientSet, err := kubernetes.NewForConfig(config)
			if err != nil {
				return errors.Errorf("Could not access Kubernetes cluster, check kubeconfig: %s", err)
			}

			toPod, err := kuberneteshelper.PodFromString(args[0])
			if err != nil {
				return errors.New("invalid DESTINATION_POD")
			}

			ingress.ToPod(clientSet, toPod)

			return nil
		},
		Example: `TODO`,
	}
	return cmd
}
