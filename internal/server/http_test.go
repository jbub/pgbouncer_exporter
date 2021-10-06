package server

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jbub/pgbouncer_exporter/internal/collector"
	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/domain"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

	"github.com/pashagolub/pgxmock"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"
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
				StoreTimeout:    time.Millisecond * 200,
			}

			conn, err := pgxmock.NewConn()
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close(context.Background())

			srv := newTestingServer(cfg, sqlstore.New(conn))
			defer srv.Close()

			if cfg.ExportPools {
				conn.ExpectQuery("SHOW POOLS").WillReturnRows(pgxmock.NewRows([]string{"database"}).AddRow("mydb"))
			}

			if cfg.ExportStats {
				conn.ExpectQuery("SHOW STATS").WillReturnRows(pgxmock.NewRows([]string{"database"}).AddRow("mydb"))
			}

			if cfg.ExportDatabases {
				conn.ExpectQuery("SHOW DATABASES").WillReturnRows(pgxmock.NewRows([]string{"database"}).AddRow("mydb"))
			}

			if cfg.ExportLists {
				conn.ExpectQuery("SHOW LISTS").WillReturnRows(pgxmock.NewRows([]string{"list"}).AddRow("mylist"))
			}

			client := srv.Client()
			resp, err := client.Get(srv.URL + cfg.TelemetryPath)
			require.NoError(t, err)
			defer resp.Body.Close()

			metrics, err := parser.TextToMetricFamilies(resp.Body)
			require.NoError(t, err)

			for _, expMetric := range testCase.metrics {
				if _, ok := metrics[expMetric]; !ok {
					require.FailNow(t, "metric not found", expMetric)
				}
			}
		})
	}
}
