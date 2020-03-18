package store

import (
	"context"

	"github.com/jbub/pgbouncer_exporter/internal/domain"
)

// NewMockStore returns new MockStore.
func NewMockStore() *MockStore {
	return &MockStore{}
}

// MockStore is a Store implementation used for testing.
type MockStore struct {
	Stats     []domain.Stat
	Pools     []domain.Pool
	Databases []domain.Database
	Lists     []domain.List
	Servers   []domain.Server
	Clients   []domain.Client

	StatsCalled     bool
	PoolsCalled     bool
	DatabasesCalled bool
	ListsCalled     bool
	ServersCalled   bool
	ClientsCalled   bool
	CloseCalled     bool
	CheckCalled     bool
}

// GetStats returns stats.
func (s *MockStore) GetStats(ctx context.Context) ([]domain.Stat, error) {
	s.StatsCalled = true
	return s.Stats, nil
}

// GetPools returns pools.
func (s *MockStore) GetPools(ctx context.Context) ([]domain.Pool, error) {
	s.PoolsCalled = true
	return s.Pools, nil
}

// GetDatabases returns databases.
func (s *MockStore) GetDatabases(ctx context.Context) ([]domain.Database, error) {
	s.DatabasesCalled = true
	return s.Databases, nil
}

// GetLists returns lists.
func (s *MockStore) GetLists(ctx context.Context) ([]domain.List, error) {
	s.ListsCalled = true
	return s.Lists, nil
}

// GetServers returns servers.
func (s *MockStore) GetServers(ctx context.Context) ([]domain.Server, error) {
	s.ServersCalled = true
	return s.Servers, nil
}

// GetClients returns clients.
func (s *MockStore) GetClients(ctx context.Context) ([]domain.Client, error) {
	s.ClientsCalled = true
	return s.Clients, nil
}

// Check checks the health of the store.
func (s *MockStore) Check(ctx context.Context) error {
	s.CheckCalled = true
	return nil
}

// Close closes the store.
func (s *MockStore) Close() error {
	s.CloseCalled = true
	return nil
}
