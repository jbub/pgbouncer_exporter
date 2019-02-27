package config

import (
	"time"

	"github.com/urfave/cli"
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
}
