package collector

import (
	"context"
	"testing"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/sqlstore"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestGetStoreResultExportEnabled(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cfg := config.Config{
		ExportStats:     true,
		ExportPools:     true,
		ExportDatabases: true,
		ExportLists:     true,
	}

	exp := New(cfg, sqlstore.New(db))
	ctx := context.Background()

	mock.ExpectQuery("SHOW STATS").WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectQuery("SHOW POOLS").WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectQuery("SHOW DATABASES").WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectQuery("SHOW LISTS").WillReturnRows(sqlmock.NewRows(nil))

	_, err = exp.getStoreResult(ctx)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStoreResultExportDisabled(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cfg := config.Config{
		ExportStats:     false,
		ExportPools:     false,
		ExportDatabases: false,
		ExportLists:     false,
	}

	exp := New(cfg, sqlstore.New(db))
	ctx := context.Background()

	_, err = exp.getStoreResult(ctx)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
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
