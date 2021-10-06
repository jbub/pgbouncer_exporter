package collector

import (
	"context"
	"testing"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

	"github.com/pashagolub/pgxmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestGetStoreResultExportEnabled(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	cfg := config.Config{
		ExportStats:     true,
		ExportPools:     true,
		ExportDatabases: true,
		ExportLists:     true,
	}

	exp := New(cfg, sqlstore.New(conn))
	ctx := context.Background()

	conn.ExpectQuery("SHOW STATS").WillReturnRows(pgxmock.NewRows(nil))
	conn.ExpectQuery("SHOW POOLS").WillReturnRows(pgxmock.NewRows(nil))
	conn.ExpectQuery("SHOW DATABASES").WillReturnRows(pgxmock.NewRows(nil))
	conn.ExpectQuery("SHOW LISTS").WillReturnRows(pgxmock.NewRows(nil))

	_, err = exp.getStoreResult(ctx)
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())
}

func TestGetStoreResultExportDisabled(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	cfg := config.Config{
		ExportStats:     false,
		ExportPools:     false,
		ExportDatabases: false,
		ExportLists:     false,
	}

	exp := New(cfg, sqlstore.New(conn))
	ctx := context.Background()

	_, err = exp.getStoreResult(ctx)
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())
}

var (
	parseLabelsCases = []struct {
		name     string
		value    string
		expected prometheus.Labels
	}{
		{
			name:     "empty",
			value:    "",
			expected: nil,
		},
		{
			name:     "invalid item",
			value:    "key",
			expected: prometheus.Labels{},
		},
		{
			name:  "blank item",
			value: "key=",
			expected: prometheus.Labels{
				"key": "",
			},
		},
		{
			name:  "single item",
			value: "key=value",
			expected: prometheus.Labels{
				"key": "value",
			},
		},
		{
			name:  "multiple items",
			value: "key=value key2=value2",
			expected: prometheus.Labels{
				"key":  "value",
				"key2": "value2",
			},
		},
		{
			name:  "multiple items with space",
			value: "key=value key2=value2 ",
			expected: prometheus.Labels{
				"key":  "value",
				"key2": "value2",
			},
		},
	}
)

func TestParseLabels(t *testing.T) {
	for _, cs := range parseLabelsCases {
		t.Run(cs.name, func(t *testing.T) {
			labels := parseLabels(cs.value)
			require.Equal(t, cs.expected, labels)
		})
	}
}
