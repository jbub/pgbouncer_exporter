package collector

import (
	"context"
	"testing"

	"github.com/jbub/pgbouncer_exporter/internal/config"
	"github.com/jbub/pgbouncer_exporter/internal/store"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestGetStoreResultExportEnabled(t *testing.T) {
	cfg := config.Config{
		ExportStats:     true,
		ExportPools:     true,
		ExportDatabases: true,
		ExportLists:     true,
	}
	st := store.NewMockStore()
	defer st.Close()

	exp := New(cfg, st)
	ctx := context.Background()

	res, err := exp.getStoreResult(ctx)
	require.NoError(t, err)

	require.True(t, st.StatsCalled)
	require.True(t, st.PoolsCalled)
	require.True(t, st.DatabasesCalled)
	require.True(t, st.ListsCalled)

	require.Equal(t, res.stats, st.Stats)
	require.Equal(t, res.pools, st.Pools)
	require.Equal(t, res.databases, st.Databases)
	require.Equal(t, res.lists, st.Lists)
}

func TestGetStoreResultExportDisabled(t *testing.T) {
	cfg := config.Config{
		ExportStats:     false,
		ExportPools:     false,
		ExportDatabases: false,
		ExportLists:     false,
	}
	st := store.NewMockStore()
	defer st.Close()

	exp := New(cfg, st)
	ctx := context.Background()

	res, err := exp.getStoreResult(ctx)
	require.NoError(t, err)

	require.False(t, st.StatsCalled)
	require.False(t, st.PoolsCalled)
	require.False(t, st.DatabasesCalled)
	require.False(t, st.ListsCalled)

	require.Equal(t, res.stats, st.Stats)
	require.Equal(t, res.pools, st.Pools)
	require.Equal(t, res.databases, st.Databases)
	require.Equal(t, res.lists, st.Lists)
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
