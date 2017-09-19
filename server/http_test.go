package server

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/jbub/pgbouncer_exporter/collector"
	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/domain"
	"github.com/jbub/pgbouncer_exporter/store"

	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

const (
	buildInfoMetric = "pgbouncer_exporter_build_info"
)

var (
	testCases = []struct {
		name            string
		exportPools     bool
		exportDatabases bool
		exportStats     bool
		exportLists     bool
		metrics         []string
	}{
		{
			name:        "stats",
			exportStats: true,
			metrics: []string{
				buildInfoMetric,
				metricName(collector.SubsystemStats, "total_requests"),
				metricName(collector.SubsystemStats, "total_received"),
				metricName(collector.SubsystemStats, "total_sent"),
				metricName(collector.SubsystemStats, "total_query_time"),
			},
		},
		{
			name:        "pools",
			exportPools: true,
			metrics: []string{
				buildInfoMetric,
				metricName(collector.SubsystemPools, "waiting_clients"),
				metricName(collector.SubsystemPools, "active_clients"),
			},
		},
		{
			name:            "databases",
			exportDatabases: true,
			metrics: []string{
				buildInfoMetric,
				metricName(collector.SubsystemDatabases, "current_connections"),
			},
		},
		{
			name:        "lists",
			exportLists: true,
			metrics: []string{
				buildInfoMetric,
				metricName(collector.SubsystemLists, "items"),
			},
		},
	}
)

func metricName(subsystem string, name string) string {
	return fmt.Sprintf("%v_%v_%v", collector.Name, subsystem, name)
}

func newTestingStore() *store.MockStore {
	st := store.NewMockStore()
	st.Stats = []domain.Stat{
		{
			Database:       "xx",
			TotalRequests:  20,
			TotalQueryTime: 344,
			TotalSent:      203,
			TotalReceived:  203,
		},
		{
			Database:       "yy",
			TotalRequests:  20,
			TotalQueryTime: 344,
			TotalSent:      203,
			TotalReceived:  203,
		},
	}
	st.Pools = []domain.Pool{
		{
			Database: "xx",
			PoolMode: "transaction",
			Active:   4,
		},
		{
			Database: "yy",
			PoolMode: "session",
			Active:   6,
		},
	}
	st.Databases = []domain.Database{
		{
			Name:               "xx",
			Database:           "xx",
			PoolMode:           "transaction",
			CurrentConnections: 4,
		},
		{
			Name:               "yy",
			Database:           "yy",
			PoolMode:           "session",
			CurrentConnections: 6,
		},
	}
	st.Lists = []domain.List{
		{
			List:  "xx",
			Items: 54,
		},
		{
			List:  "yy",
			Items: 68,
		},
	}
	return st
}

func newTestingServer(cfg config.Config, st domain.Store) *httptest.Server {
	exp := collector.New(cfg, st)
	httpSrv := New(cfg, exp)
	return httptest.NewServer(httpSrv.srv.Handler)
}

func TestResponseContainsMetrics(t *testing.T) {
	var parser expfmt.TextParser

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := config.Config{
				TelemetryPath:   "/metrics",
				ExportPools:     testCase.exportPools,
				ExportDatabases: testCase.exportDatabases,
				ExportStats:     testCase.exportStats,
				ExportLists:     testCase.exportLists,
			}

			st := newTestingStore()
			defer st.Close()

			srv := newTestingServer(cfg, st)
			defer srv.Close()

			client := srv.Client()
			resp, err := client.Get(srv.URL + cfg.TelemetryPath)
			assert.NoError(t, err)
			defer resp.Body.Close()

			metrics, err := parser.TextToMetricFamilies(resp.Body)
			assert.NoError(t, err)

			if cfg.ExportPools {
				assert.True(t, st.PoolsCalled)
			}
			if cfg.ExportStats {
				assert.True(t, st.StatsCalled)
			}
			if cfg.ExportDatabases {
				assert.True(t, st.DatabasesCalled)
			}
			if cfg.ExportLists {
				assert.True(t, st.ListsCalled)
			}

			for _, expMetric := range testCase.metrics {
				if _, ok := metrics[expMetric]; !ok {
					assert.FailNow(t, "metric not found", expMetric)
				}
			}
		})
	}
}
