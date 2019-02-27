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
)

type gaugeVecItem struct {
	gaugeVec *prometheus.GaugeVec
	resolve  resolveGaugeVecFunc
	enabled  bool
}

type gaugeVecValueItem struct {
	labels []string
	count  float64
}

type resolveGaugeVecFunc func(res *storeResult) []gaugeVecValueItem

type storeResult struct {
	stats     []domain.Stat
	pools     []domain.Pool
	databases []domain.Database
	lists     []domain.List
}

// Exporter represents pgbouncer prometheus stats exporter.
type Exporter struct {
	cfg           config.Config
	st            domain.Store
	mutex         sync.RWMutex // guards Collect
	gaugeVecItems []gaugeVecItem
}

// New returns new Exporter.
func New(cfg config.Config, st domain.Store) *Exporter {
	return &Exporter{
		st:  st,
		cfg: cfg,
		gaugeVecItems: []gaugeVecItem{
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_requests",
					Help:      "Total number of SQL requests pooled by pgbouncer.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalRequests),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_received",
					Help:      "Total number of SQL requests pooled by pgbouncer.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalReceived),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_sent",
					Help:      "Total volume in bytes of network traffic sent by pgbouncer.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalSent),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_query_time",
					Help:      "Total number of microseconds spent by pgbouncer when actively connected to PostgreSQL.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalQueryTime),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_xact_time",
					Help:      "Total number of microseconds spent by pgbouncer when connected to PostgreSQL in a transaction, either idle in transaction or executing queries.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalXactTime),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_query_count",
					Help:      "Total number of SQL queries pooled by pgbouncer.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalQueryCount),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportStats,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemStats,
					Name:      "total_xact_count",
					Help:      "Total number of SQL transactions pooled by pgbouncer.",
				}, []string{"database"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, stat := range res.stats {
						items = append(items, gaugeVecValueItem{
							labels: []string{stat.Database},
							count:  float64(stat.TotalXactCount),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "active_clients",
					Help:      "Client connections that are linked to server connection and can process queries.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.Active),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "waiting_clients",
					Help:      "Client connections have sent queries but have not yet got a server connection.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.Waiting),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "active_server",
					Help:      "Server connections that linked to client.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.ServerActive),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "idle_server",
					Help:      "Server connections that unused and immediately usable for client queries.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.ServerIdle),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "used_server",
					Help:      "Server connections that have been idle more than server_check_delay, so they needs server_check_query to run on it before it can be used.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.ServerUsed),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "tested_server",
					Help:      "Server connections that are currently running either server_reset_query or server_check_query.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.ServerTested),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "login_server",
					Help:      "Server connections currently in logging in process.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.ServerLogin),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportPools,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemPools,
					Name:      "max_wait",
					Help:      "How long the first (oldest) client in queue has waited, in seconds. If this starts increasing, then the current pool of servers does not handle requests quick enough. Reason may be either overloaded server or just too small of a pool_size setting.",
				}, []string{"database", "user", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, pool := range res.pools {
						items = append(items, gaugeVecValueItem{
							labels: []string{pool.Database, pool.User, pool.PoolMode},
							count:  float64(pool.MaxWait),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportDatabases,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemDatabases,
					Name:      "current_connections",
					Help:      "Current number of connections.",
				}, []string{"name", "pool_mode"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, database := range res.databases {
						items = append(items, gaugeVecValueItem{
							labels: []string{database.Name, database.PoolMode},
							count:  float64(database.CurrentConnections),
						})
					}
					return items
				},
			},
			{
				enabled: cfg.ExportLists,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemLists,
					Name:      "items",
					Help:      "List of internal pgbouncer information.",
				}, []string{"list"}),
				resolve: func(res *storeResult) (items []gaugeVecValueItem) {
					for _, list := range res.lists {
						items = append(items, gaugeVecValueItem{
							labels: []string{list.List},
							count:  float64(list.Items),
						})
					}
					return items
				},
			},
		},
	}
}

// Describe implements prometheus Collector.Describe.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, gaugeVecItem := range e.gaugeVecItems {
		if gaugeVecItem.enabled {
			gaugeVecItem.gaugeVec.Describe(ch)
		}
	}
}

// Collect implements prometheus Collector.Collect.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.StoreTimeout)
	defer cancel()

	res, err := e.getStoreResult(ctx)
	if err != nil {
		log.Errorf("unable to get store result: %v", err)
		return
	}

	for _, gaugeVecItem := range e.gaugeVecItems {
		if !gaugeVecItem.enabled {
			continue
		}

		items := gaugeVecItem.resolve(res)
		for _, item := range items {
			gaugeVecItem.gaugeVec.WithLabelValues(item.labels...).Set(item.count)
		}

		gaugeVecItem.gaugeVec.Collect(ch)
	}
}

func (e *Exporter) getStoreResult(ctx context.Context) (*storeResult, error) {
	res := new(storeResult)

	if e.cfg.ExportStats {
		stats, err := e.st.GetStats(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get stats: %v", err)
		}
		res.stats = stats
	}

	if e.cfg.ExportPools {
		pools, err := e.st.GetPools(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get pools: %v", err)
		}
		res.pools = pools
	}

	if e.cfg.ExportDatabases {
		databases, err := e.st.GetDatabases(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get databases: %v", err)
		}
		res.databases = databases
	}

	if e.cfg.ExportLists {
		lists, err := e.st.GetLists(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get lists: %v", err)
		}
		res.lists = lists
	}

	return res, nil
}