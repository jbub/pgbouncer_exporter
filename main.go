package main

import (
	"log"
	"os"
	"time"

	"github.com/jbub/pgbouncer_exporter/cmd"
	"github.com/jbub/pgbouncer_exporter/internal/collector"

	"github.com/prometheus/common/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  collector.Name,
		Usage: collector.Name,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "web.listen-address",
				Usage:   "Address on which to expose metrics and web interface.",
				EnvVars: []string{"WEB_LISTEN_ADDRESS"},
				Value:   ":9127",
			},
			&cli.StringFlag{
				Name:    "web.telemetry-path",
				Usage:   "Path under which to expose metrics.",
				EnvVars: []string{"WEB_TELEMETRY_PATH"},
				Value:   "/metrics",
			},
			&cli.StringFlag{
				Name:    "database-url",
				Usage:   "Database connection url.",
				EnvVars: []string{"DATABASE_URL"},
			},
			&cli.BoolFlag{
				Name:    "export-stats",
				Usage:   "Export stats.",
				EnvVars: []string{"EXPORT_STATS"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "export-pools",
				Usage:   "Export pools.",
				EnvVars: []string{"EXPORT_POOLS"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "export-databases",
				Usage:   "Export databases.",
				EnvVars: []string{"EXPORT_DATABASES"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "export-lists",
				Usage:   "Export lists.",
				EnvVars: []string{"EXPORT_LISTS"},
				Value:   true,
			},
			&cli.DurationFlag{
				Name:    "store-timeout",
				Usage:   "Per method store timeout.",
				EnvVars: []string{"STORE_TIMEOUT"},
				Value:   time.Second * 2,
			},
			&cli.StringFlag{
				Name:    "default-labels",
				Usage:   "Default prometheus labels applied to all metrics. Format: label1=value1 label2=value2",
				EnvVars: []string{"DEFAULT_LABELS"},
			},
		},
		Commands: []*cli.Command{
			cmd.Server,
			cmd.Health,
		},
		Version: version.Info(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
