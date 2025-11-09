package debug

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
)

// JobRepositoryImpl is a repository implementation that loads job data from a JSON file
type JobRepositoryImpl struct {
	data *Data
}

// NewJobRepository creates a new debug job repository from loaded data
func NewJobRepository(data *Data) repository.JobRepository {
	return &JobRepositoryImpl{
		data: data,
	}
}

func (j *JobRepositoryImpl) FetchActiveJobs(_ context.Context, _, _, _ string) ([]*entity.Job, error) {
	return j.data.Jobs, nil
}
