package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jbub/pgbouncer_exporter/internal/collector"
	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/server"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/common/version"
	"github.com/urfave/cli/v2"
)

// Server is a cli command used for running exporter http server.
var Server = &cli.Command{
	Name:   "server",
	Usage:  "Starts exporter server.",
	Action: runServer,
}

func runServer(ctx *cli.Context) error {
	cfg := config.LoadFromCLI(ctx)

	connCfg, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("could not parse connection config: %v", err)
	}
	connCfg.PreferSimpleProtocol = true

	conn, err := pgx.ConnectConfig(context.Background(), connCfg)
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}
	defer conn.Close(context.Background())

	store := sqlstore.New(conn)

	checkCtx, cancel := context.WithTimeout(context.Background(), cfg.StoreTimeout)
	defer cancel()

	if err := store.Check(checkCtx); err != nil {
		return fmt.Errorf("could not check store: %v", err)
	}

	exp := collector.New(cfg, store)
	srv := server.New(cfg, exp)

	log.Println("Starting ", collector.Name, version.Info())
	log.Println("Server listening on", cfg.ListenAddress)
	log.Println("Metrics available at", cfg.TelemetryPath)
	log.Println("Build context", version.BuildContext())

	if err := srv.Run(); err != nil {
		return fmt.Errorf("could not run server: %v", err)
	}
	return nil
}
