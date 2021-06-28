package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/openservicemesh/osm-trouble/pkg/logger"
)

var	log = logger.New("osm-trouble")


func main() {
	app := &cli.App{
		Name:  "OSM Troubleshooter",
		Usage: "osm-trouble <area>",
		Action: func(c *cli.Context) error {
			fmt.Println("Testing")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
