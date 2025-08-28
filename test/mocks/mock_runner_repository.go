package mocks

import (
	"context"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

type MockRunnerRepository struct {
	ListRunnersFunc    func(ctx context.Context, owner, repo string) ([]*entities.Runner, error)
	ListOrgRunnersFunc func(ctx context.Context, org string) ([]*entities.Runner, error)
	GetRunnerFunc      func(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error)
}

func (m *MockRunnerRepository) ListRunners(ctx context.Context, owner, repo string) ([]*entities.Runner, error) {
	if m.ListRunnersFunc != nil {
		return m.ListRunnersFunc(ctx, owner, repo)
	}
	return []*entities.Runner{}, nil
}

func (m *MockRunnerRepository) ListOrgRunners(ctx context.Context, org string) ([]*entities.Runner, error) {
	if m.ListOrgRunnersFunc != nil {
		return m.ListOrgRunnersFunc(ctx, org)
	}
	return []*entities.Runner{}, nil
}

func (m *MockRunnerRepository) GetRunner(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error) {
	if m.GetRunnerFunc != nil {
		return m.GetRunnerFunc(ctx, owner, repo, runnerID)
	}
	return nil, nil
}

func NewMockRunnerRepository() *MockRunnerRepository {
	return &MockRunnerRepository{}
}

func (m *MockRunnerRepository) WithListRunners(f func(ctx context.Context, owner, repo string) ([]*entities.Runner, error)) *MockRunnerRepository {
	m.ListRunnersFunc = f
	return m
}

func (m *MockRunnerRepository) WithListOrgRunners(f func(ctx context.Context, org string) ([]*entities.Runner, error)) *MockRunnerRepository {
	m.ListOrgRunnersFunc = f
	return m
}

func (m *MockRunnerRepository) WithGetRunner(f func(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error)) *MockRunnerRepository {
	m.GetRunnerFunc = f
	return m
}