package domain

import (
	"context"
)

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

type List struct {
	List  string
	Items int64
}

type Store interface {
	GetStats(ctx context.Context) ([]Stat, error)

	GetPools(ctx context.Context) ([]Pool, error)

	GetDatabases(ctx context.Context) ([]Database, error)

	GetLists(ctx context.Context) ([]List, error)

	Close()
}
