package domain

import (
	"context"
)

// Stat represents stat row.
type Stat struct {
	Database        string
	TotalRequests   int64
	TotalReceived   int64
	TotalSent       int64
	TotalQueryTime  int64
	AverageRequests int64
	AverageReceived int64
	AverageSent     int64
	AverageQuery    int64
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
	Check() error

	// Close closes the store.
	Close() error
}
