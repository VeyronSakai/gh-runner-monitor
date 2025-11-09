package debug

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
)

// RunnerRepositoryImpl is a repository implementation that loads data from a JSON file
type RunnerRepositoryImpl struct {
	data *Data
}

// NewRunnerRepository creates a new debug repository from loaded data
func NewRunnerRepository(data *Data) repository.RunnerRepository {
	return &RunnerRepositoryImpl{
		data: data,
	}
}



func (d *RunnerRepositoryImpl) FetchRunners(_ context.Context, _, _, _ string) ([]*entity.Runner, error) {
	return d.data.Runners, nil
}
