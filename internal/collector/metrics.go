package collector

import (
	"github.com/jbub/pgbouncer_exporter/internal/config"

	"github.com/prometheus/client_golang/prometheus"
)

func buildMetrics(cfg config.Config) []metric {
	return []metric{
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_requests"),
			help:    "Total number of SQL requests pooled by pgbouncer.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalRequests),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_received"),
			help:    "Total volume in bytes of network traffic received by pgbouncer.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalReceived),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_sent"),
			help:    "Total volume in bytes of network traffic sent by pgbouncer.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalSent),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_query_time"),
			help:    "Total number of microseconds spent by pgbouncer when actively connected to PostgreSQL.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalQueryTime),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_xact_time"),
			help:    "Total number of microseconds spent by pgbouncer when connected to PostgreSQL in a transaction, either idle in transaction or executing queries.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalXactTime),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_query_count"),
			help:    "Total number of SQL queries pooled by pgbouncer.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalQueryCount),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportStats,
			name:    fqName(SubsystemStats, "total_xact_count"),
			help:    "Total number of SQL transactions pooled by pgbouncer.",
			labels:  []string{"database"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, stat := range res.stats {
					results = append(results, metricResult{
						labels: []string{stat.Database},
						value:  float64(stat.TotalXactCount),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "active_clients"),
			help:    "Client connections that are linked to server connection and can process queries.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.Active),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "waiting_clients"),
			help:    "Client connections have sent queries but have not yet got a server connection.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.Waiting),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "active_server"),
			help:    "Server connections that are linked to a client.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.ServerActive),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "idle_server"),
			help:    "Server connections that are unused and immediately usable for client queries.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.ServerIdle),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "used_server"),
			help:    "Server connections that have been idle for more than server_check_delay, so they need server_check_query to run on them before they can be used again.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.ServerUsed),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "tested_server"),
			help:    "Server connections that are currently running either server_reset_query or server_check_query.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.ServerTested),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "login_server"),
			help:    "Server connections currently in the process of logging in.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.ServerLogin),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportPools,
			name:    fqName(SubsystemPools, "max_wait"),
			help:    "How long the first (oldest) client in the queue has waited, in seconds. If this starts increasing, then the current pool of servers does not handle requests quickly enough. The reason may be either an overloaded server or just too small of a pool_size setting.",
			labels:  []string{"database", "user", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, pool := range res.pools {
					results = append(results, metricResult{
						labels: []string{pool.Database, pool.User, pool.PoolMode},
						value:  float64(pool.MaxWait),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportDatabases,
			name:    fqName(SubsystemDatabases, "current_connections"),
			help:    "Current number of connections for this database.",
			labels:  []string{"name", "pool_mode"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, database := range res.databases {
					results = append(results, metricResult{
						labels: []string{database.Name, database.PoolMode},
						value:  float64(database.CurrentConnections),
					})
				}
				return results
			},
		},
		{
			enabled: cfg.ExportLists,
			name:    fqName(SubsystemLists, "items"),
			help:    "List of internal pgbouncer information.",
			labels:  []string{"list"},
			valType: prometheus.GaugeValue,
			eval: func(res *storeResult) (results []metricResult) {
				for _, list := range res.lists {
					results = append(results, metricResult{
						labels: []string{list.List},
						value:  float64(list.Items),
					})
				}
				return results
			},
		},
	}
}

func fqName(subsystem string, name string) string {
	return prometheus.BuildFQName(Name, subsystem, name)
}
