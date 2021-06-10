package main

import (
	"errors"
	"log"
	"os"

	"github.com/larkintuckerllc/balancer/internal/balancer"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "balancer",
		Usage: "exports metrics",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "cluster",
				Usage:    "cluster names",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "idle",
				Usage:    "idle (in minutes) - greater than or equal to 0",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "hpa",
				Usage:    "hpa name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "location",
				Usage:    "cluster location",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "namespace",
				Usage:    "hpa namespace",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "prefix",
				Usage:    "context prefix",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "project",
				Usage:    "project ID",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "value",
				Usage:    "metric value - greater than or equal to 100",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			clusters := c.StringSlice("cluster")
			hpa := c.String("hpa")
			idle := c.Int("idle")
			location := c.String("location")
			namespace := c.String("namespace")
			prefix := c.String("prefix")
			project := c.String("project")
			value := c.Int("value")
			if idle < 0 {
				err := errors.New("idle must be greater than or equal to 0")
				return err
			}
			if value < 100 {
				err := errors.New("value must be greater than or equal to 100")
				return err
			}
			err := balancer.Execute(project, location, prefix, clusters, namespace, hpa, value, idle)
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
