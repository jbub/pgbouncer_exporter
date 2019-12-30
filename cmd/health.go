package cmd

import (
	"context"
	"fmt"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/store"

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
	st, err := store.NewSQL(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to initialize store: %v", err)
	}
	defer st.Close()

	checkCtx, cancel := context.WithTimeout(context.Background(), cfg.StoreTimeout)
	defer cancel()

	if err := st.Check(checkCtx); err != nil {
		return fmt.Errorf("store health check failed: %v", err)
	}
	return nil
}
