package main

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/openservicemesh/osm-health/pkg/cli"
	"github.com/openservicemesh/osm-health/pkg/connectivity"
	"github.com/openservicemesh/osm-health/pkg/kubernetes/pod"
)

const connectivityPodToURLDesc = `
Checks connectivity between a Kubernetes pod and a host name (or URL)
`

const connectivityPodToURLExample = `$ osm-health connectivity pod-to-url source-namespace/source-pod https://contoso.com/store`

func newConnectivityPodToURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pod-to-url source-namespace/source-pod destination-url",
		Short:   "Checks connectivity between a Kubernetes pod and a given URL",
		Example: connectivityPodToURLExample,
		Long:    connectivityPodToURLDesc,
		Args:    cli.ExactArgsWithError(2, errors.New("requires 2 arguments: source-namespace/source-pod destination-url")),
		RunE: func(_ *cobra.Command, args []string) error {
			srcPod, err := pod.FromString(args[0])
			if err != nil {
				return errors.New("invalid source-namespace/source-pod")
			}

			dstURL, err := url.Parse(args[1])
			if err != nil {
				return errors.New("invalid destination-url")
			}

			connectivity.PodToURL(srcPod, dstURL, settings.Namespace())
			return nil
		},
	}
}
