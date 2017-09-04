package config

import (
	"time"
)

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
