package store

import (
	"context"
	"database/sql"

	"github.com/jbub/pgbouncer_exporter/internal/domain"

	"github.com/jmoiron/sqlx"
)

type stat struct {
	Database          string `db:"database"`
	TotalRequests     int64  `db:"total_requests"`
	TotalReceived     int64  `db:"total_received"`
	TotalSent         int64  `db:"total_sent"`
	TotalQueryTime    int64  `db:"total_query_time"`
	TotalXactCount    int64  `db:"total_xact_count"`
	TotalXactTime     int64  `db:"total_xact_time"`
	TotalQueryCount   int64  `db:"total_query_count"`
	TotalWaitTime     int64  `db:"total_wait_time"`
	AverageRequests   int64  `db:"avg_req"`
	AverageReceived   int64  `db:"avg_recv"`
	AverageSent       int64  `db:"avg_sent"`
	AverageQuery      int64  `db:"avg_query"`
	AverageQueryCount int64  `db:"avg_query_count"`
	AverageQueryTime  int64  `db:"avg_query_time"`
	AverageXactTime   int64  `db:"avg_xact_time"`
	AverageXactCount  int64  `db:"avg_xact_count"`
	AverageWaitTime   int64  `db:"avg_wait_time"`
}

type pool struct {
	Database            string         `db:"database"`
	User                string         `db:"user"`
	Active              int64          `db:"cl_active"`
	Waiting             int64          `db:"cl_waiting"`
	ActiveCancelReq     int64          `db:"cl_active_cancel_req"`
	WaitingCancelReq    int64          `db:"cl_waiting_cancel_req"`
	ServerActive        int64          `db:"sv_active"`
	ServerActiveCancel  int64          `db:"sv_active_cancel"`
	ServerBeingCanceled int64          `db:"sv_being_canceled"`
	ServerIdle          int64          `db:"sv_idle"`
	ServerUsed          int64          `db:"sv_used"`
	ServerTested        int64          `db:"sv_tested"`
	ServerLogin         int64          `db:"sv_login"`
	MaxWait             int64          `db:"maxwait"`
	MaxWaitUs           int64          `db:"maxwait_us"`
	PoolMode            sql.NullString `db:"pool_mode"`
}

type database struct {
	Name               string         `db:"name"`
	Host               sql.NullString `db:"host"`
	Port               int64          `db:"port"`
	Database           string         `db:"database"`
	ForceUser          sql.NullString `db:"force_user"`
	PoolSize           int64          `db:"pool_size"`
	ReservePool        int64          `db:"reserve_pool"`
	PoolMode           sql.NullString `db:"pool_mode"`
	MaxConnections     int64          `db:"max_connections"`
	CurrentConnections int64          `db:"current_connections"`
	Paused             int64          `db:"paused"`
	Disabled           int64          `db:"disabled"`
}

type list struct {
	List  string `db:"list"`
	Items int64  `db:"items"`
}

type server struct {
	Type            string         `db:"type"`
	User            string         `db:"user"`
	Database        string         `db:"database"`
	State           string         `db:"state"`
	Addr            string         `db:"addr"`
	Port            int64          `db:"port"`
	LocalAddr       string         `db:"local_addr"`
	LocalPort       int64          `db:"local_port"`
	ConnectTime     string         `db:"connect_time"`
	RequestTime     string         `db:"request_time"`
	Wait            int64          `db:"wait"`
	WaitUs          int64          `db:"wait_us"`
	CloseNeeded     int64          `db:"close_needed"`
	Ptr             string         `db:"ptr"`
	Link            string         `db:"link"`
	RemotePid       string         `db:"remote_pid"`
	Tls             string         `db:"tls"`
	ApplicationName sql.NullString `db:"application_name"`
}

type client struct {
	Type            string         `db:"type"`
	User            string         `db:"user"`
	Database        string         `db:"database"`
	State           string         `db:"state"`
	Addr            string         `db:"addr"`
	Port            int64          `db:"port"`
	LocalAddr       string         `db:"local_addr"`
	LocalPort       int64          `db:"local_port"`
	ConnectTime     string         `db:"connect_time"`
	RequestTime     string         `db:"request_time"`
	Wait            int64          `db:"wait"`
	WaitUs          int64          `db:"wait_us"`
	CloseNeeded     int64          `db:"close_needed"`
	Ptr             string         `db:"ptr"`
	Link            string         `db:"link"`
	RemotePid       string         `db:"remote_pid"`
	Tls             string         `db:"tls"`
	ApplicationName sql.NullString `db:"application_name"`
}

// NewSQL returns a new SQLStore.
func NewSQL(dataSource string) (*SQLStore, error) {
	db, err := sqlx.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	return &SQLStore{db: db}, nil
}

// SQLStore is a sql based Store implementation.
type SQLStore struct {
	db *sqlx.DB
}

// GetStats returns stats.
func (s *SQLStore) GetStats(ctx context.Context) ([]domain.Stat, error) {
	var stats []stat
	if err := s.db.SelectContext(ctx, &stats, "SHOW STATS"); err != nil {
		return nil, err
	}
	var result []domain.Stat
	for _, row := range stats {
		result = append(result, domain.Stat(row))
	}
	return result, nil
}

// GetPools returns pools.
func (s *SQLStore) GetPools(ctx context.Context) ([]domain.Pool, error) {
	var pools []pool
	if err := s.db.SelectContext(ctx, &pools, "SHOW POOLS"); err != nil {
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
func (s *SQLStore) GetDatabases(ctx context.Context) ([]domain.Database, error) {
	var databases []database
	if err := s.db.SelectContext(ctx, &databases, "SHOW DATABASES"); err != nil {
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
func (s *SQLStore) GetLists(ctx context.Context) ([]domain.List, error) {
	var lists []list
	if err := s.db.SelectContext(ctx, &lists, "SHOW LISTS"); err != nil {
		return nil, err
	}
	var result []domain.List
	for _, row := range lists {
		result = append(result, domain.List(row))
	}
	return result, nil
}

// GetServers returns servers.
func (s *SQLStore) GetServers(ctx context.Context) ([]domain.Server, error) {
	var servers []server
	if err := s.db.SelectContext(ctx, &servers, "SHOW SERVERS"); err != nil {
		return nil, err
	}
	var result []domain.Server
	for _, row := range servers {
		result = append(result, domain.Server{
			Database:        row.Database,
			State:           row.State,
			ApplicationName: row.ApplicationName.String,
		})
	}
	return result, nil
}

// GetClients returns servers.
func (s *SQLStore) GetClients(ctx context.Context) ([]domain.Client, error) {
	var clients []client
	if err := s.db.SelectContext(ctx, &clients, "SHOW CLIENTS"); err != nil {
		return nil, err
	}
	var result []domain.Client
	for _, row := range clients {
		result = append(result, domain.Client{
			Database:        row.Database,
			State:           row.State,
			ApplicationName: row.ApplicationName.String,
		})
	}
	return result, nil
}

// Check checks the health of the store.
func (s *SQLStore) Check(ctx context.Context) error {
	// we cant use db.Ping because it is making a ";" sql query which pgbouncer does not support
	rows, err := s.db.QueryContext(ctx, "SHOW VERSION")
	if err != nil {
		return err
	}
	return rows.Close()
}

// Close closes the store.
func (s *SQLStore) Close() error {
	return s.db.Close()
}
