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
		DefaultLabels:   ctx.String("default-labels"),
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
	DefaultLabels   string
}
