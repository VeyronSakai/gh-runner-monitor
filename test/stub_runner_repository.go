package test

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// StubRunnerRepository is a stub implementation of repository.RunnerRepository for testing.
// It returns pre-configured responses without any behavior verification.
type StubRunnerRepository struct {
	// Runners is the data that will be returned by GetRunners
	Runners []*entity.Runner
	// GetRunnersError is the error that will be returned by GetRunners
	GetRunnersError error
}

// NewStubRunnerRepositoryWithError creates a new StubRunnerRepository that returns errors.
func NewStubRunnerRepositoryWithError(getRunnersErr error) *StubRunnerRepository {
	return &StubRunnerRepository{
		GetRunnersError: getRunnersErr,
	}
}

func (s *StubRunnerRepository) FetchRunners(_ context.Context, _, _, _ string) ([]*entity.Runner, error) {
	if s.GetRunnersError != nil {
		return nil, s.GetRunnersError
	}
	return s.Runners, nil
}
