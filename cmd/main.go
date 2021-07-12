package main

import (
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/openservicemesh/osm-trouble/pkg/connectivity"
	"github.com/openservicemesh/osm-trouble/pkg/kubernetes"
	"github.com/openservicemesh/osm-trouble/pkg/logger"
)

var log = logger.New("osm-trouble/main")

func main() {
	app := &cli.App{
		Name:  "OSM Troubleshooter",
		Usage: "osm-trouble <area>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "from-pod",
				Aliases: []string{"fp"},
				Usage:   "Test connectivity from Pod",
			},
			&cli.StringFlag{
				Name:    "to-pod",
				Aliases: []string{"tp"},
				Usage:   "Test connectivity to Pod",
			},
		},
		Action: func(c *cli.Context) error {
			verb := c.Args().First()

			if verb == "connectivity" {
				fromPodString := c.String("from-pod")
				if fromPodString == "" {
					log.Fatal().Msg("--from-pod is required")
				}

				toPodString := c.String("to-pod")
				if toPodString == "" {
					log.Fatal().Msg("--to-pod is required")
				}

				fromPod, err := kubernetes.PodFromString(fromPodString)
				if err != nil {
					log.Fatal().Msg("Err fetching Pod ")
				}

				toPod, err := kubernetes.PodFromString(toPodString)
				if err != nil {
					log.Fatal().Msg("Err fetching Pod ")
				}

				connectivity.PodToPod(fromPod, toPod)
				return nil
			}

			if verb == "collect" {
				return nil
			}

			log.Fatal().Msgf("What's the verb? I don't recognize %q", verb)

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
