package main

import (
	goflag "flag"
	"io"
	"os"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"

	"github.com/openservicemesh/osm-health/pkg/cli"
	"github.com/openservicemesh/osm-health/pkg/logger"
)

var globalUsage = `The osm-health cli enables you to
	(1) check osm health status
	(2) debug osm issues
`

var settings = cli.New()

var log = logger.New("osm-health/main")

func newRootCmd(config *action.Configuration, in io.Reader, out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "osm-health",
		Short:        "Check Open Service Mesh health status and debug issues",
		Long:         globalUsage,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)

	// Add subcommands here
	cmd.AddCommand(
		newCollectCmd(config, in, out),
		newConnectivityCmd(config, in, out),
		newValidateCmd(),
	)

	_ = flags.Parse(args)

	return cmd
}

func initCommands() *cobra.Command {
	actionConfig := new(action.Configuration)
	cmd := newRootCmd(actionConfig, os.Stdin, os.Stdout, os.Args[1:])
	_ = actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "secret", debug)

	// run when each command's execute method is called
	cobra.OnInitialize(func() {
		if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "secret", debug); err != nil {
			os.Exit(1)
		}
	})

	return cmd
}

func main() {
	log.Info().Msg(`👋🏻`)
	cmd := initCommands()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func debug(format string, v ...interface{}) {
}
