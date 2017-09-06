package collector

import (
	"context"
	"testing"

	"github.com/jbub/pgbouncer_exporter/config"
	"github.com/jbub/pgbouncer_exporter/store"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	assert.True(t, st.StatsCalled)
	assert.True(t, st.PoolsCalled)
	assert.True(t, st.DatabasesCalled)
	assert.True(t, st.ListsCalled)

	assert.Equal(t, res.stats, st.Stats)
	assert.Equal(t, res.pools, st.Pools)
	assert.Equal(t, res.databases, st.Databases)
	assert.Equal(t, res.lists, st.Lists)
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
	assert.NoError(t, err)

	assert.False(t, st.StatsCalled)
	assert.False(t, st.PoolsCalled)
	assert.False(t, st.DatabasesCalled)
	assert.False(t, st.ListsCalled)

	assert.Equal(t, res.stats, st.Stats)
	assert.Equal(t, res.pools, st.Pools)
	assert.Equal(t, res.databases, st.Databases)
	assert.Equal(t, res.lists, st.Lists)
}
