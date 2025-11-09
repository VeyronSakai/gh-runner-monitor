package test

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// StubJobRepository is a stub implementation of repository.JobRepository for testing.
type StubJobRepository struct {
	// Jobs is the data that will be returned by GetActiveJobs
	Jobs []*entity.Job
	// GetActiveJobsError is the error that will be returned by GetActiveJobs
	GetActiveJobsError error
}

func (s *StubJobRepository) FetchActiveJobs(_ context.Context, _, _, _ string) ([]*entity.Job, error) {
	if s.GetActiveJobsError != nil {
		return nil, s.GetActiveJobsError
	}
	return s.Jobs, nil
}
