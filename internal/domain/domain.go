package domain

import (
	"context"
)

// Stat represents stat row.
type Stat struct {
	Database                     string
	TotalReceived                int64
	TotalSent                    int64
	TotalQueryTime               int64
	TotalXactCount               int64
	TotalXactTime                int64
	TotalQueryCount              int64
	TotalWaitTime                int64
	TotalServerAssignmentCount   int64
	AverageReceived              int64
	AverageSent                  int64
	AverageQueryCount            int64
	AverageQueryTime             int64
	AverageXactTime              int64
	AverageXactCount             int64
	AverageWaitTime              int64
	AverageServerAssignmentCount int64
}

// Pool represents pool row.
type Pool struct {
	Database            string
	User                string
	Active              int64
	Waiting             int64
	CancelReq           int64
	ActiveCancelReq     int64
	WaitingCancelReq    int64
	ServerActive        int64
	ServerActiveCancel  int64
	ServerBeingCanceled int64
	ServerIdle          int64
	ServerUsed          int64
	ServerTested        int64
	ServerLogin         int64
	MaxWait             int64
	MaxWaitUs           int64
	PoolMode            string
}

// Database represents database row.
type Database struct {
	Name               string
	Host               string
	Port               int64
	Database           string
	ForceUser          string
	PoolSize           int64
	MinPoolSize        int64
	ReservePool        int64
	PoolMode           string
	MaxConnections     int64
	CurrentConnections int64
	Paused             int64
	Disabled           int64
}

// List represents list row.
type List struct {
	List  string
	Items int64
}

// Store defines interface for accessing pgbouncer stats.
type Store interface {
	// GetStats returns stats.
	GetStats(ctx context.Context) ([]Stat, error)

	// GetPools returns pools.
	GetPools(ctx context.Context) ([]Pool, error)

	// GetDatabases returns databases.
	GetDatabases(ctx context.Context) ([]Database, error)

	// GetLists returns lists.
	GetLists(ctx context.Context) ([]List, error)

	// Check checks the health of the store.
	Check(ctx context.Context) error
}
