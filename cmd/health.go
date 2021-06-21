package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

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

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("could not open db: %v", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	checkCtx, cancel := context.WithTimeout(context.Background(), cfg.StoreTimeout)
	defer cancel()

	if err := store.Check(checkCtx); err != nil {
		return fmt.Errorf("store health check failed: %v", err)
	}
	return nil
}
