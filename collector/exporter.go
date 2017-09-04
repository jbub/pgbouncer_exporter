package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	Name              = "pgbouncer_exporter"
	SubsystemStats    = "stats"
	SubsystemPools    = "pools"
	SubsystemDatabase = "database"
	SubsystemLists    = "lists"
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
				enabled: cfg.ExportDatabases,
				gaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: Name,
					Subsystem: SubsystemDatabase,
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

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, gaugeVecItem := range e.gaugeVecItems {
		if gaugeVecItem.enabled {
			gaugeVecItem.gaugeVec.Describe(ch)
		}
	}
}

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
