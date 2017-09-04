package store

import (
	"context"

	"github.com/jbub/pgbouncer_exporter/domain"
)

func NewMockStore() *MockStore {
	return &MockStore{}
}

type MockStore struct {
	Stats     []domain.Stat
	Pools     []domain.Pool
	Databases []domain.Database
	Lists     []domain.List

	StatsCalled     bool
	PoolsCalled     bool
	DatabasesCalled bool
	ListsCalled     bool
	CloseCalled     bool
}

func (s *MockStore) GetStats(ctx context.Context) ([]domain.Stat, error) {
	s.StatsCalled = true
	return s.Stats, nil
}

func (s *MockStore) GetPools(ctx context.Context) ([]domain.Pool, error) {
	s.PoolsCalled = true
	return s.Pools, nil
}

func (s *MockStore) GetDatabases(ctx context.Context) ([]domain.Database, error) {
	s.DatabasesCalled = true
	return s.Databases, nil
}

func (s *MockStore) GetLists(ctx context.Context) ([]domain.List, error) {
	s.ListsCalled = true
	return s.Lists, nil
}

func (s *MockStore) Close() {
	s.CloseCalled = true
}
