package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/domain"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// Name is the name of the exporter.
	Name = "pgbouncer_exporter"
)

// Names of the exporter subsystems.
const (
	SubsystemStats     = "stats"
	SubsystemPools     = "pools"
	SubsystemDatabases = "database"
	SubsystemLists     = "lists"
	SubsystemServers   = "servers"
	SubsystemClients   = "clients"
)

var (
	_ prometheus.Collector = &Exporter{}
)

type metric struct {
	enabled bool
	desc    *prometheus.Desc
	valType prometheus.ValueType
	eval    func(res *storeResult) []metricResult
}

type metricResult struct {
	labels []string
	value  float64
}

type storeResult struct {
	stats     []domain.Stat
	pools     []domain.Pool
	databases []domain.Database
	lists     []domain.List
	servers   []domain.Server
	clients   []domain.Client
}

// Exporter represents pgbouncer prometheus stats exporter.
type Exporter struct {
	cfg     config.Config
	stor    domain.Store
	mut     sync.Mutex // guards Collect
	metrics []metric
}

// New returns new Exporter.
func New(cfg config.Config, stor domain.Store) *Exporter {
	return &Exporter{
		stor: stor,
		cfg:  cfg,
		metrics: []metric{
			{
				enabled: cfg.ExportStats,
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_requests"),
					"Total number of SQL requests pooled by pgbouncer.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_received"),
					"Total volume in bytes of network traffic received by pgbouncer.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_sent"),
					"Total volume in bytes of network traffic sent by pgbouncer.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_query_time"),
					"Total number of microseconds spent by pgbouncer when actively connected to PostgreSQL.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_xact_time"),
					"Total number of microseconds spent by pgbouncer when connected to PostgreSQL in a transaction, either idle in transaction or executing queries.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_query_count"),
					"Total number of SQL queries pooled by pgbouncer.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemStats, "total_xact_count"),
					"Total number of SQL transactions pooled by pgbouncer.",
					[]string{"database"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "active_clients"),
					"Client connections that are linked to server connection and can process queries.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "waiting_clients"),
					"Client connections have sent queries but have not yet got a server connection.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "active_server"),
					"Server connections that are linked to a client.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "idle_server"),
					"Server connections that are unused and immediately usable for client queries.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "used_server"),
					"Server connections that have been idle for more than server_check_delay, so they need server_check_query to run on them before they can be used again.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "tested_server"),
					"Server connections that are currently running either server_reset_query or server_check_query.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "login_server"),
					"Server connections currently in the process of logging in.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemPools, "max_wait"),
					"How long the first (oldest) client in the queue has waited, in seconds. If this starts increasing, then the current pool of servers does not handle requests quickly enough. The reason may be either an overloaded server or just too small of a pool_size setting.",
					[]string{"database", "user", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemDatabases, "current_connections"),
					"Current number of connections for this database.",
					[]string{"name", "pool_mode"},
					nil,
				),
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
				desc: prometheus.NewDesc(
					fqName(SubsystemLists, "items"),
					"List of internal pgbouncer information.",
					[]string{"list"},
					nil,
				),
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
			{
				enabled: cfg.ExportServers,
				desc: prometheus.NewDesc(
					fqName(SubsystemServers, "active"),
					"Active servers.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findServerResults(res.servers, "active")
				},
			},
			{
				enabled: cfg.ExportServers,
				desc: prometheus.NewDesc(
					fqName(SubsystemServers, "used"),
					"Used servers.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findServerResults(res.servers, "used")
				},
			},
			{
				enabled: cfg.ExportServers,
				desc: prometheus.NewDesc(
					fqName(SubsystemServers, "idle"),
					"Idle servers.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findServerResults(res.servers, "idle")
				},
			},
			{
				enabled: cfg.ExportClients,
				desc: prometheus.NewDesc(
					fqName(SubsystemClients, "active"),
					"Active clients.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findClientResults(res.clients, "active")
				},
			},
			{
				enabled: cfg.ExportClients,
				desc: prometheus.NewDesc(
					fqName(SubsystemClients, "used"),
					"Used clients.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findClientResults(res.clients, "used")
				},
			},
			{
				enabled: cfg.ExportClients,
				desc: prometheus.NewDesc(
					fqName(SubsystemClients, "waiting"),
					"Waiting clients.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findClientResults(res.clients, "waiting")
				},
			},
			{
				enabled: cfg.ExportClients,
				desc: prometheus.NewDesc(
					fqName(SubsystemClients, "idle"),
					"Idle clients.",
					[]string{"database", "user"},
					nil,
				),
				valType: prometheus.GaugeValue,
				eval: func(res *storeResult) (results []metricResult) {
					return findClientResults(res.clients, "idle")
				},
			},
		},
	}
}

// Describe implements prometheus Collector.Describe.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, met := range e.metrics {
		if !met.enabled {
			continue
		}
		ch <- met.desc
	}
}

// Collect implements prometheus Collector.Collect.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mut.Lock()
	defer e.mut.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.StoreTimeout)
	defer cancel()

	res, err := e.getStoreResult(ctx)
	if err != nil {
		log.Errorf("could not get store result: %v", err)
		return
	}

	for _, met := range e.metrics {
		if !met.enabled {
			continue
		}

		results := met.eval(res)

		for _, res := range results {
			ch <- prometheus.MustNewConstMetric(
				met.desc,
				met.valType,
				res.value,
				res.labels...,
			)
		}
	}
}

func (e *Exporter) getStoreResult(ctx context.Context) (*storeResult, error) {
	res := new(storeResult)

	if e.cfg.ExportStats {
		stats, err := e.stor.GetStats(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get stats: %v", err)
		}
		res.stats = stats
	}

	if e.cfg.ExportPools {
		pools, err := e.stor.GetPools(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get pools: %v", err)
		}
		res.pools = pools
	}

	if e.cfg.ExportDatabases {
		databases, err := e.stor.GetDatabases(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get databases: %v", err)
		}
		res.databases = databases
	}

	if e.cfg.ExportLists {
		lists, err := e.stor.GetLists(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get lists: %v", err)
		}
		res.lists = lists
	}

	if e.cfg.ExportServers {
		servers, err := e.stor.GetServers(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get servers: %v", err)
		}
		res.servers = servers
	}

	if e.cfg.ExportClients {
		clients, err := e.stor.GetClients(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get clients: %v", err)
		}
		res.clients = clients
	}

	return res, nil
}

func fqName(subsystem string, name string) string {
	return prometheus.BuildFQName(Name, subsystem, name)
}

func findServerResults(servers []domain.Server, state string) []metricResult {
	var results = []metricResult{}
	for _, server := range servers {
		if server.State == state {
			results = append(results, metricResult{
				labels: []string{server.Database, server.User},
				value:  1.0,
			})
		}
	}
	return results
}

func findClientResults(clients []domain.Client, state string) []metricResult {
	var results = []metricResult{}
	for _, client := range clients {
		if client.State == state {
			results = append(results, metricResult{
				labels: []string{client.Database, client.User},
				value:  1.0,
			})
		}
	}
	return results
}
