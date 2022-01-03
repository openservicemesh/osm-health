package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/cli"
	"github.com/openservicemesh/osm-health/pkg/ingress"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
)

const ingressToPodExample = `$ osm-health ingress to-pod destination-namespace/destination-pod`

func newIngressToPodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "to-pod destination-namespace/destination-pod",
		Short:   "Checks ingress to a given Kubernetes pod",
		Example: ingressToPodExample,
		Long:    `Checks ingress to a given Kubernetes pod`,
		Args:    cli.ExactArgsWithError(1, errors.New("requires 1 argument: destination-namespace/destination-pod")),
		RunE: func(_ *cobra.Command, args []string) error {
			log.Info().Msgf("Checking Ingress to Pod %s", args[0])

			client, err := pod.GetKubeClient()
			if err != nil {
				return err
			}

			dstPod, err := pod.FromString(args[0])
			if err != nil {
				return errors.New("invalid destination-namespace/destination-pod")
			}

			osmControlPlaneNamespace := settings.Namespace()

			ingress.ToDestinationPod(client, dstPod, osmControlPlaneNamespace)

			return nil
		},
	}
	return cmd
}
