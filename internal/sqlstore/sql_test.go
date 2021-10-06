package sqlstore

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestGetStats(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	st := New(conn)
	data := map[string]interface{}{
		"database":          "pgbouncer",
		"total_requests":    int64(1),
		"total_received":    int64(2),
		"total_sent":        int64(3),
		"total_query_time":  int64(4),
		"total_xact_count":  int64(5),
		"total_xact_time":   int64(6),
		"total_query_count": int64(7),
		"total_wait_time":   int64(8),
		"avg_req":           int64(9),
		"avg_recv":          int64(10),
		"avg_sent":          int64(11),
		"avg_query":         int64(12),
		"avg_query_count":   int64(13),
		"avg_query_time":    int64(14),
		"avg_xact_time":     int64(15),
		"avg_xact_count":    int64(16),
		"avg_wait_time":     int64(17),
	}

	conn.ExpectQuery("SHOW STATS").WillReturnRows(mapToRows(data))

	stats, err := st.GetStats(context.Background())
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())

	stat := stats[0]
	require.Equal(t, data["database"], stat.Database)
	require.Equal(t, data["total_requests"], stat.TotalRequests)
	require.Equal(t, data["total_received"], stat.TotalReceived)
	require.Equal(t, data["total_sent"], stat.TotalSent)
	require.Equal(t, data["total_query_time"], stat.TotalQueryTime)
	require.Equal(t, data["total_xact_count"], stat.TotalXactCount)
	require.Equal(t, data["total_xact_time"], stat.TotalXactTime)
	require.Equal(t, data["total_query_count"], stat.TotalQueryCount)
	require.Equal(t, data["total_wait_time"], stat.TotalWaitTime)
	require.Equal(t, data["avg_req"], stat.AverageRequests)
	require.Equal(t, data["avg_recv"], stat.AverageReceived)
	require.Equal(t, data["avg_sent"], stat.AverageSent)
	require.Equal(t, data["avg_query"], stat.AverageQuery)
	require.Equal(t, data["avg_query_count"], stat.AverageQueryCount)
	require.Equal(t, data["avg_query_time"], stat.AverageQueryTime)
	require.Equal(t, data["avg_xact_time"], stat.AverageXactTime)
	require.Equal(t, data["avg_xact_count"], stat.AverageXactCount)
	require.Equal(t, data["avg_wait_time"], stat.AverageWaitTime)
}

func TestGetPools(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	st := New(conn)

	data := map[string]interface{}{
		"database":   "pgbouncer",
		"user":       "myuser",
		"cl_active":  int64(1),
		"cl_waiting": int64(2),
		"sv_active":  int64(3),
		"sv_idle":    int64(4),
		"sv_used":    int64(5),
		"sv_tested":  int64(6),
		"sv_login":   int64(7),
		"maxwait":    int64(8),
		"maxwait_us": int64(9),
		"pool_mode":  "transaction",
	}

	conn.ExpectQuery("SHOW POOLS").WillReturnRows(mapToRows(data))

	pools, err := st.GetPools(context.Background())
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())

	pool := pools[0]
	require.Equal(t, data["database"], pool.Database)
	require.Equal(t, data["user"], pool.User)
	require.Equal(t, data["cl_active"], pool.Active)
	require.Equal(t, data["cl_waiting"], pool.Waiting)
	require.Equal(t, data["sv_active"], pool.ServerActive)
	require.Equal(t, data["sv_idle"], pool.ServerIdle)
	require.Equal(t, data["sv_used"], pool.ServerUsed)
	require.Equal(t, data["sv_tested"], pool.ServerTested)
	require.Equal(t, data["sv_login"], pool.ServerLogin)
	require.Equal(t, data["maxwait"], pool.MaxWait)
	require.Equal(t, data["maxwait_us"], pool.MaxWaitUs)
	require.Equal(t, data["pool_mode"], pool.PoolMode)
}

func TestGetDatabases(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	st := New(conn)

	data := map[string]interface{}{
		"database":            "pgbouncer",
		"name":                "myname",
		"host":                "localhost",
		"port":                int64(23),
		"force_user":          "myuser",
		"pool_size":           int64(4),
		"reserve_pool":        int64(5),
		"pool_mode":           "transaction",
		"max_connections":     int64(7),
		"current_connections": int64(8),
		"paused":              int64(9),
		"disabled":            int64(10),
	}

	conn.ExpectQuery("SHOW DATABASES").WillReturnRows(mapToRows(data))

	databases, err := st.GetDatabases(context.Background())
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())

	database := databases[0]
	require.Equal(t, data["database"], database.Database)
	require.Equal(t, data["name"], database.Name)
	require.Equal(t, data["host"], database.Host)
	require.Equal(t, data["port"], database.Port)
	require.Equal(t, data["force_user"], database.ForceUser)
	require.Equal(t, data["pool_size"], database.PoolSize)
	require.Equal(t, data["reserve_pool"], database.ReservePool)
	require.Equal(t, data["pool_mode"], database.PoolMode)
	require.Equal(t, data["max_connections"], database.MaxConnections)
	require.Equal(t, data["current_connections"], database.CurrentConnections)
	require.Equal(t, data["paused"], database.Paused)
	require.Equal(t, data["disabled"], database.Disabled)
}

func TestGetLists(t *testing.T) {
	conn, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(context.Background())

	st := New(conn)

	data := map[string]interface{}{
		"list":  "mylist",
		"items": int64(6),
	}

	conn.ExpectQuery("SHOW LISTS").WillReturnRows(mapToRows(data))

	lists, err := st.GetLists(context.Background())
	require.NoError(t, err)
	require.NoError(t, conn.ExpectationsWereMet())

	list := lists[0]
	require.Equal(t, data["list"], list.List)
	require.Equal(t, data["items"], list.Items)
}

func mapToRows(data map[string]interface{}) *pgxmock.Rows {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}
	rows := pgxmock.NewRows(columns)
	rows.AddRow(values...)
	return rows
}
