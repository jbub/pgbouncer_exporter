package config

import (
	"time"

	"github.com/urfave/cli/v2"
)

func LoadFromCLI(ctx *cli.Context) Config {
	return Config{
		ListenAddress:   ctx.String("web.listen-address"),
		TelemetryPath:   ctx.String("web.telemetry-path"),
		DatabaseURL:     ctx.String("database-url"),
		StoreTimeout:    ctx.Duration("store-timeout"),
		ExportStats:     ctx.Bool("export-stats"),
		ExportPools:     ctx.Bool("export-pools"),
		ExportDatabases: ctx.Bool("export-databases"),
		ExportLists:     ctx.Bool("export-lists"),
		ExportServers:   ctx.Bool("export-servers"),
		ExportClients:   ctx.Bool("export-clients"),
	}
}

// Config represents exporter configuration.
type Config struct {
	ListenAddress string
	TelemetryPath string
	DatabaseURL   string
	StoreTimeout  time.Duration

	ExportStats     bool
	ExportPools     bool
	ExportDatabases bool
	ExportLists     bool
	ExportServers   bool
	ExportClients   bool
}
