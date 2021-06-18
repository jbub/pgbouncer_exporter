package collector

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/domain"

	"github.com/prometheus/client_golang/prometheus"
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

var (
	_ prometheus.Collector = &Exporter{}
)

type metric struct {
	enabled bool
	name    string
	help    string
	labels  []string
	valType prometheus.ValueType
	eval    func(res *storeResult) []metricResult
}

func (m metric) desc(constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(m.name, m.help, m.labels, constLabels)
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
}

// Exporter represents pgbouncer prometheus stats exporter.
type Exporter struct {
	cfg         config.Config
	stor        domain.Store
	mut         sync.Mutex // guards Collect
	constLabels prometheus.Labels
	metrics     []metric
}

// New returns new Exporter.
func New(cfg config.Config, stor domain.Store) *Exporter {
	return &Exporter{
		stor:        stor,
		cfg:         cfg,
		constLabels: parseLabels(cfg.DefaultLabels),
		metrics:     buildMetrics(cfg),
	}
}

// Describe implements prometheus Collector.Describe.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, met := range e.metrics {
		if !met.enabled {
			continue
		}
		ch <- met.desc(e.constLabels)
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
		log.Printf("could not get store result: %v", err)
		return
	}

	for _, met := range e.metrics {
		if !met.enabled {
			continue
		}

		results := met.eval(res)

		for _, res := range results {
			ch <- prometheus.MustNewConstMetric(
				met.desc(e.constLabels),
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
			return nil, fmt.Errorf("could not get stats: %v", err)
		}
		res.stats = stats
	}

	if e.cfg.ExportPools {
		pools, err := e.stor.GetPools(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get pools: %v", err)
		}
		res.pools = pools
	}

	if e.cfg.ExportDatabases {
		databases, err := e.stor.GetDatabases(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get databases: %v", err)
		}
		res.databases = databases
	}

	if e.cfg.ExportLists {
		lists, err := e.stor.GetLists(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get lists: %v", err)
		}
		res.lists = lists
	}

	return res, nil
}

func parseLabels(s string) prometheus.Labels {
	if s == "" {
		return nil
	}

	items := strings.Split(s, " ")
	res := make(prometheus.Labels, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if parts := strings.SplitN(item, "=", 2); len(parts) == 2 {
			res[parts[0]] = parts[1]
		}
	}
	return res
}
