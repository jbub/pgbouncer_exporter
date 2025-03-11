package sqlstore

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	st := New(db)

	data := map[string]any{
		"database":                      "pgbouncer",
		"total_xact_count":              1,
		"total_query_count":             2,
		"total_received":                3,
		"total_sent":                    4,
		"total_xact_time":               5,
		"total_query_time":              6,
		"total_wait_time":               7,
		"total_server_assignment_count": 8,
		"avg_xact_count":                9,
		"avg_query_count":               10,
		"avg_recv":                      11,
		"avg_sent":                      12,
		"avg_xact_time":                 13,
		"avg_query_time":                14,
		"avg_wait_time":                 15,
		"avg_server_assignment_count":   16,
	}

	mock.ExpectQuery("SHOW STATS").WillReturnRows(mapToRows(data))

	stats, err := st.GetStats(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	stat := stats[0]
	require.Equal(t, data["database"].(string), stat.Database)
	require.Equal(t, int64(data["total_xact_count"].(int)), stat.TotalXactCount)
	require.Equal(t, int64(data["total_query_count"].(int)), stat.TotalQueryCount)
	require.Equal(t, int64(data["total_received"].(int)), stat.TotalReceived)
	require.Equal(t, int64(data["total_sent"].(int)), stat.TotalSent)
	require.Equal(t, int64(data["total_xact_time"].(int)), stat.TotalXactTime)
	require.Equal(t, int64(data["total_query_time"].(int)), stat.TotalQueryTime)
	require.Equal(t, int64(data["total_wait_time"].(int)), stat.TotalWaitTime)
	require.Equal(t, int64(data["total_server_assignment_count"].(int)), stat.TotalServerAssignmentCount)
	require.Equal(t, int64(data["avg_xact_count"].(int)), stat.AverageXactCount)
	require.Equal(t, int64(data["avg_query_count"].(int)), stat.AverageQueryCount)
	require.Equal(t, int64(data["avg_recv"].(int)), stat.AverageReceived)
	require.Equal(t, int64(data["avg_sent"].(int)), stat.AverageSent)
	require.Equal(t, int64(data["avg_xact_time"].(int)), stat.AverageXactTime)
	require.Equal(t, int64(data["avg_query_time"].(int)), stat.AverageQueryTime)
	require.Equal(t, int64(data["avg_wait_time"].(int)), stat.AverageWaitTime)
	require.Equal(t, int64(data["avg_server_assignment_count"].(int)), stat.AverageServerAssignmentCount)
}

func TestGetPools(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	st := New(db)

	data := map[string]any{
		"database":   "pgbouncer",
		"user":       "myuser",
		"cl_active":  1,
		"cl_waiting": 2,
		"sv_active":  3,
		"sv_idle":    4,
		"sv_used":    5,
		"sv_tested":  6,
		"sv_login":   7,
		"maxwait":    8,
		"maxwait_us": 9,
		"pool_mode":  "transaction",
	}

	mock.ExpectQuery("SHOW POOLS").WillReturnRows(mapToRows(data))

	pools, err := st.GetPools(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	pool := pools[0]
	require.Equal(t, data["database"].(string), pool.Database)
	require.Equal(t, data["user"].(string), pool.User)
	require.Equal(t, int64(data["cl_active"].(int)), pool.Active)
	require.Equal(t, int64(data["cl_waiting"].(int)), pool.Waiting)
	require.Equal(t, int64(data["sv_active"].(int)), pool.ServerActive)
	require.Equal(t, int64(data["sv_idle"].(int)), pool.ServerIdle)
	require.Equal(t, int64(data["sv_used"].(int)), pool.ServerUsed)
	require.Equal(t, int64(data["sv_tested"].(int)), pool.ServerTested)
	require.Equal(t, int64(data["sv_login"].(int)), pool.ServerLogin)
	require.Equal(t, int64(data["maxwait"].(int)), pool.MaxWait)
	require.Equal(t, int64(data["maxwait_us"].(int)), pool.MaxWaitUs)
	require.Equal(t, data["pool_mode"].(string), pool.PoolMode)
}

func TestGetDatabases(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	st := New(db)

	data := map[string]any{
		"database":            "pgbouncer",
		"name":                "myname",
		"host":                "localhost",
		"port":                23,
		"force_user":          "myuser",
		"pool_size":           4,
		"reserve_pool_size":   5,
		"reserve_pool":        5,
		"pool_mode":           "transaction",
		"max_connections":     7,
		"current_connections": 8,
		"paused":              9,
		"disabled":            10,
		"server_lifetime":     11,
	}

	mock.ExpectQuery("SHOW DATABASES").WillReturnRows(mapToRows(data))

	databases, err := st.GetDatabases(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	database := databases[0]
	require.Equal(t, data["database"].(string), database.Database)
	require.Equal(t, data["name"].(string), database.Name)
	require.Equal(t, data["host"].(string), database.Host)
	require.Equal(t, int64(data["port"].(int)), database.Port)
	require.Equal(t, data["force_user"].(string), database.ForceUser)
	require.Equal(t, int64(data["pool_size"].(int)), database.PoolSize)
	require.Equal(t, int64(data["reserve_pool_size"].(int)), database.ReservePoolSize)
	require.Equal(t, int64(data["reserve_pool"].(int)), database.ReservePoolSize)
	require.Equal(t, data["pool_mode"].(string), database.PoolMode)
	require.Equal(t, int64(data["max_connections"].(int)), database.MaxConnections)
	require.Equal(t, int64(data["current_connections"].(int)), database.CurrentConnections)
	require.Equal(t, int64(data["paused"].(int)), database.Paused)
	require.Equal(t, int64(data["disabled"].(int)), database.Disabled)
	require.Equal(t, int64(data["server_lifetime"].(int)), database.ServerLifetime)
}

func TestGetLists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	st := New(db)

	data := map[string]any{
		"list":  "mylist",
		"items": 6,
	}

	mock.ExpectQuery("SHOW LISTS").WillReturnRows(mapToRows(data))

	lists, err := st.GetLists(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	list := lists[0]
	require.Equal(t, data["list"].(string), list.List)
	require.Equal(t, int64(data["items"].(int)), list.Items)
}

func mapToRows(data map[string]any) *sqlmock.Rows {
	columns := make([]string, 0, len(data))
	values := make([]driver.Value, 0, len(data))
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(values...)
	return rows
}
