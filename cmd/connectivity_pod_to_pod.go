package main

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// osmConfigClient "github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned"
	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kubernetesHelper"
	smiAccessClient "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/access/clientset/versioned"
)

const connectivityPodToPodDesc = `
Checks connectivity between Kubernetes pods
	(add more descriptive description)
`

const connectivityPodToPodExample = `
Example:
	(add example)
`

type connectivityPodToPodCmd struct {
	out             io.Writer
	srcPod          string
	destPod         string
	clientSet       kubernetes.Interface
	smiAccessClient smiAccessClient.Interface
	// meshConfigClient osmConfigClient.Interface
	restConfig *rest.Config
}

func newConnectivityPodToPodCmd(config *action.Configuration, in io.Reader, out io.Writer) *cobra.Command {
	podToPodCmd := &connectivityPodToPodCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "pod-to-pod SOURCE_POD DESTINATION_POD",
		Short: "Checks connectivity between Kubernetes pods",
		Long:  connectivityPodToPodDesc,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			podToPodCmd.srcPod = args[0]
			podToPodCmd.destPod = args[1]

			config, err := settings.RESTClientGetter().ToRESTConfig()
			if err != nil {
				return errors.Errorf("Error fetching kubeconfig: %s", err)
			}

			podToPodCmd.restConfig = config

			clientSet, err := kubernetes.NewForConfig(config)
			if err != nil {
				return errors.Errorf("Could not access Kubernetes cluster, check kubeconfig: %s", err)
			}
			podToPodCmd.clientSet = clientSet

			accessClient, err := smiAccessClient.NewForConfig(config)
			if err != nil {
				return errors.Errorf("Could not initialize SMI Access client: %s", err)
			}
			podToPodCmd.smiAccessClient = accessClient

			// configClient, err := osmConfigClient.NewForConfig(config)
			// if err != nil {
			// 	return errors.Errorf("Could not initialize OSM Config client: %s", err)
			// }
			// podToPodCmd.meshConfigClient = configClient

			return podToPodCmd.run()
		},
		Example: connectivityPodToPodExample,
	}
	return cmd
}

func (podToPodCmd *connectivityPodToPodCmd) run() error {
	srcPod, err := kubernetesHelper.PodFromString(podToPodCmd.srcPod)
	if err != nil {
		return errors.New("invalid SOURCE_POD")
	}

	destPod, err := kubernetesHelper.PodFromString(podToPodCmd.destPod)
	if err != nil {
		return errors.New("invalid DESTINATION_POD")
	}

	connectivity.PodToPod(srcPod, destPod, podToPodCmd.clientSet)
	return nil
}
