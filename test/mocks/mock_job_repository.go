package mocks

import (
	"context"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

type MockJobRepository struct {
	ListActiveJobsFunc    func(ctx context.Context, owner, repo string) ([]*entities.Job, error)
	ListOrgActiveJobsFunc func(ctx context.Context, org string) ([]*entities.Job, error)
	GetJobFunc            func(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error)
}

func (m *MockJobRepository) ListActiveJobs(ctx context.Context, owner, repo string) ([]*entities.Job, error) {
	if m.ListActiveJobsFunc != nil {
		return m.ListActiveJobsFunc(ctx, owner, repo)
	}
	return []*entities.Job{}, nil
}

func (m *MockJobRepository) ListOrgActiveJobs(ctx context.Context, org string) ([]*entities.Job, error) {
	if m.ListOrgActiveJobsFunc != nil {
		return m.ListOrgActiveJobsFunc(ctx, org)
	}
	return []*entities.Job{}, nil
}

func (m *MockJobRepository) GetJob(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error) {
	if m.GetJobFunc != nil {
		return m.GetJobFunc(ctx, owner, repo, jobID)
	}
	return nil, nil
}

func NewMockJobRepository() *MockJobRepository {
	return &MockJobRepository{}
}

func (m *MockJobRepository) WithListActiveJobs(f func(ctx context.Context, owner, repo string) ([]*entities.Job, error)) *MockJobRepository {
	m.ListActiveJobsFunc = f
	return m
}

func (m *MockJobRepository) WithListOrgActiveJobs(f func(ctx context.Context, org string) ([]*entities.Job, error)) *MockJobRepository {
	m.ListOrgActiveJobsFunc = f
	return m
}

func (m *MockJobRepository) WithGetJob(f func(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error)) *MockJobRepository {
	m.GetJobFunc = f
	return m
}