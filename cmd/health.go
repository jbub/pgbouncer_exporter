package cmd

import (
	"context"
	"fmt"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

// Health is a cli command used for checking the health of the system.
var Health = &cli.Command{
	Name:   "health",
	Usage:  "Checks the health of the system.",
	Action: checkHealth,
}

func checkHealth(ctx *cli.Context) error {
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
		return fmt.Errorf("store health check failed: %v", err)
	}
	return nil
}
