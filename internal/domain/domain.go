package domain

import (
	"context"
)

// Stat represents stat row.
type Stat struct {
	Database          string
	TotalRequests     int64
	TotalReceived     int64
	TotalSent         int64
	TotalQueryTime    int64
	TotalXactCount    int64
	TotalXactTime     int64
	TotalQueryCount   int64
	TotalWaitTime     int64
	AverageRequests   int64
	AverageReceived   int64
	AverageSent       int64
	AverageQuery      int64
	AverageQueryCount int64
	AverageQueryTime  int64
	AverageXactTime   int64
	AverageXactCount  int64
	AverageWaitTime   int64
}

// Pool represents pool row.
type Pool struct {
	Database     string
	User         string
	Active       int64
	Waiting      int64
	ServerActive int64
	ServerIdle   int64
	ServerUsed   int64
	ServerTested int64
	ServerLogin  int64
	MaxWait      int64
	MaxWaitUs    int64
	PoolMode     string
}

// Database represents database row.
type Database struct {
	Name               string
	Host               string
	Port               int64
	Database           string
	ForceUser          string
	PoolSize           int64
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

// Server represents server row.
type Server struct {
	Database        string
	State           string
	ApplicationName string
}

// Client represents client row.
type Client struct {
	Database        string
	State           string
	ApplicationName string
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

	// GetServers returns servers.
	GetServers(ctx context.Context) ([]Server, error)

	// GetClients returns clients.
	GetClients(ctx context.Context) ([]Client, error)

	// Check checks the health of the store.
	Check(ctx context.Context) error

	// Close closes the store.
	Close() error
}
