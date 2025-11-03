package test

import (
	"context"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// StubRunnerRepository is a stub implementation of repository.RunnerRepository for testing.
// It returns pre-configured responses without any behavior verification.
type StubRunnerRepository struct {
	// Runners is the data that will be returned by GetRunners
	Runners []*entity.Runner
	// Jobs is the data that will be returned by GetActiveJobs
	Jobs []*entity.Job
	// GetRunnersError is the error that will be returned by GetRunners
	GetRunnersError error
	// GetActiveJobsError is the error that will be returned by GetActiveJobs
	GetActiveJobsError error
}

// NewStubRunnerRepositoryWithError creates a new StubRunnerRepository that returns errors.
func NewStubRunnerRepositoryWithError(getRunnersErr, getActiveJobsErr error) *StubRunnerRepository {
	return &StubRunnerRepository{
		GetRunnersError:    getRunnersErr,
		GetActiveJobsError: getActiveJobsErr,
	}
}

// GetRunners returns the configured runners or error.
// This is a stub - it simply returns pre-configured data without any behavior verification.
func (s *StubRunnerRepository) GetRunners(_ context.Context, _, _, _ string) ([]*entity.Runner, error) {
	if s.GetRunnersError != nil {
		return nil, s.GetRunnersError
	}
	return s.Runners, nil
}

// GetActiveJobs returns the configured jobs or error.
// This is a stub - it simply returns pre-configured data without any behavior verification.
func (s *StubRunnerRepository) GetActiveJobs(_ context.Context, _, _, _ string) ([]*entity.Job, error) {
	if s.GetActiveJobsError != nil {
		return nil, s.GetActiveJobsError
	}
	return s.Jobs, nil
}

// GetCurrentTime returns the current time for testing
func (s *StubRunnerRepository) GetCurrentTime() time.Time {
	return time.Now()
}
