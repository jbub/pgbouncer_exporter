package cmd

import (
	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/server"

	"github.com/urfave/cli"
)

var Server = &cli.Command{
	Name:   "server",
	Usage:  "Starts exporter server.",
	Action: runServer,
}

func runServer(ctx *cli.Context) error {
	cfg := config.Config{
		ListenAddress:   ctx.String("web.listen-address"),
		TelemetryPath:   ctx.String("web.telemetry-path"),
		DatabaseURL:     ctx.String("database-url"),
		StoreTimeout:    ctx.Duration("store-timeout"),
		ExportStats:     ctx.Bool("export-stats"),
		ExportPools:     ctx.Bool("export-pools"),
		ExportDatabases: ctx.Bool("export-databases"),
		ExportLists:     ctx.Bool("export-lists"),
	}

	s, err := server.New(cfg)
	if err != nil {
		return err
	}
	return s.Run()
}
