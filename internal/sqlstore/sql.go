package sqlstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jbub/pgbouncer_exporter/internal/domain"
)

type pool struct {
	Database     string
	User         string
	Active       int64
	Waiting      int64
	CancelReq    int64
	ServerActive int64
	ServerIdle   int64
	ServerUsed   int64
	ServerTested int64
	ServerLogin  int64
	MaxWait      int64
	MaxWaitUs    int64
	PoolMode     sql.NullString
}

type database struct {
	Name               string
	Host               sql.NullString
	Port               int64
	Database           string
	ForceUser          sql.NullString
	PoolSize           int64
	MinPoolSize        int64
	ReservePool        int64
	PoolMode           sql.NullString
	MaxConnections     int64
	CurrentConnections int64
	Paused             int64
	Disabled           int64
}

// New returns a new SQLStore.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// Store is a sql based Store implementation.
type Store struct {
	db *sql.DB
}

// GetStats returns stats.
func (s *Store) GetStats(ctx context.Context) ([]domain.Stat, error) {
	rows, err := s.db.QueryContext(ctx, "SHOW STATS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var row domain.Stat
	var stats []domain.Stat

	for rows.Next() {
		dest := make([]interface{}, 0, len(columns))

		for _, column := range columns {
			switch column {
			case "database":
				dest = append(dest, &row.Database)
			case "total_requests":
				dest = append(dest, &row.TotalRequests)
			case "total_received":
				dest = append(dest, &row.TotalReceived)
			case "total_sent":
				dest = append(dest, &row.TotalSent)
			case "total_query_time":
				dest = append(dest, &row.TotalQueryTime)
			case "total_xact_count":
				dest = append(dest, &row.TotalXactCount)
			case "total_xact_time":
				dest = append(dest, &row.TotalXactTime)
			case "total_query_count":
				dest = append(dest, &row.TotalQueryCount)
			case "total_wait_time":
				dest = append(dest, &row.TotalWaitTime)
			case "avg_req":
				dest = append(dest, &row.AverageRequests)
			case "avg_recv":
				dest = append(dest, &row.AverageReceived)
			case "avg_sent":
				dest = append(dest, &row.AverageSent)
			case "avg_query":
				dest = append(dest, &row.AverageQuery)
			case "avg_query_count":
				dest = append(dest, &row.AverageQueryCount)
			case "avg_query_time":
				dest = append(dest, &row.AverageQueryTime)
			case "avg_xact_time":
				dest = append(dest, &row.AverageXactTime)
			case "avg_xact_count":
				dest = append(dest, &row.AverageXactCount)
			case "avg_wait_time":
				dest = append(dest, &row.AverageWaitTime)
			default:
				return nil, fmt.Errorf("unexpected column: %v", column)
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		stats = append(stats, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetPools returns pools.
func (s *Store) GetPools(ctx context.Context) ([]domain.Pool, error) {
	rows, err := s.db.QueryContext(ctx, "SHOW POOLS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var row pool
	var pools []pool

	for rows.Next() {
		dest := make([]interface{}, 0, len(columns))

		for _, column := range columns {
			switch column {
			case "database":
				dest = append(dest, &row.Database)
			case "user":
				dest = append(dest, &row.User)
			case "cl_active":
				dest = append(dest, &row.Active)
			case "cl_waiting":
				dest = append(dest, &row.Waiting)
			case "cl_cancel_req":
				dest = append(dest, &row.CancelReq)
			case "sv_active":
				dest = append(dest, &row.ServerActive)
			case "sv_idle":
				dest = append(dest, &row.ServerIdle)
			case "sv_used":
				dest = append(dest, &row.ServerUsed)
			case "sv_tested":
				dest = append(dest, &row.ServerTested)
			case "sv_login":
				dest = append(dest, &row.ServerLogin)
			case "maxwait":
				dest = append(dest, &row.MaxWait)
			case "maxwait_us":
				dest = append(dest, &row.MaxWaitUs)
			case "pool_mode":
				dest = append(dest, &row.PoolMode)
			default:
				return nil, fmt.Errorf("unexpected column: %v", column)
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		pools = append(pools, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []domain.Pool

	for _, row := range pools {
		result = append(result, domain.Pool{
			Database:     row.Database,
			User:         row.User,
			Active:       row.Active,
			Waiting:      row.Waiting,
			ServerActive: row.ServerActive,
			ServerIdle:   row.ServerIdle,
			ServerUsed:   row.ServerUsed,
			ServerTested: row.ServerTested,
			ServerLogin:  row.ServerLogin,
			MaxWait:      row.MaxWait,
			MaxWaitUs:    row.MaxWaitUs,
			PoolMode:     row.PoolMode.String,
		})
	}

	return result, nil
}

// GetDatabases returns databases.
func (s *Store) GetDatabases(ctx context.Context) ([]domain.Database, error) {
	rows, err := s.db.QueryContext(ctx, "SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var row database
	var databases []database

	for rows.Next() {
		dest := make([]interface{}, 0, len(columns))

		for _, column := range columns {
			switch column {
			case "database":
				dest = append(dest, &row.Database)
			case "name":
				dest = append(dest, &row.Name)
			case "host":
				dest = append(dest, &row.Host)
			case "port":
				dest = append(dest, &row.Port)
			case "force_user":
				dest = append(dest, &row.ForceUser)
			case "pool_size":
				dest = append(dest, &row.PoolSize)
			case "min_pool_size":
				dest = append(dest, &row.MinPoolSize)
			case "reserve_pool":
				dest = append(dest, &row.ReservePool)
			case "pool_mode":
				dest = append(dest, &row.PoolMode)
			case "max_connections":
				dest = append(dest, &row.MaxConnections)
			case "current_connections":
				dest = append(dest, &row.CurrentConnections)
			case "paused":
				dest = append(dest, &row.Paused)
			case "disabled":
				dest = append(dest, &row.Disabled)
			default:
				return nil, fmt.Errorf("unexpected column: %v", column)
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		databases = append(databases, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []domain.Database

	for _, row := range databases {
		result = append(result, domain.Database{
			Name:               row.Name,
			Host:               row.Host.String,
			Port:               row.Port,
			Database:           row.Database,
			ForceUser:          row.ForceUser.String,
			PoolSize:           row.PoolSize,
			ReservePool:        row.ReservePool,
			PoolMode:           row.PoolMode.String,
			MaxConnections:     row.MaxConnections,
			CurrentConnections: row.CurrentConnections,
			Paused:             row.Paused,
			Disabled:           row.Disabled,
		})
	}

	return result, nil
}

// GetLists returns lists.
func (s *Store) GetLists(ctx context.Context) ([]domain.List, error) {
	rows, err := s.db.QueryContext(ctx, "SHOW LISTS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var row domain.List
	var lists []domain.List

	for rows.Next() {
		dest := make([]interface{}, 0, len(columns))

		for _, column := range columns {
			switch column {
			case "list":
				dest = append(dest, &row.List)
			case "items":
				dest = append(dest, &row.Items)
			default:
				return nil, fmt.Errorf("unexpected column: %v", column)
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		lists = append(lists, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

// Check checks the health of the store.
func (s *Store) Check(ctx context.Context) error {
	// we cant use db.Ping because it is making a ";" sql query which pgbouncer does not support
	rows, err := s.db.QueryContext(ctx, "SHOW VERSION")
	if err != nil {
		return err
	}
	return rows.Close()
}
