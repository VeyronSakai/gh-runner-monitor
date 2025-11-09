package repository

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// JobRepository defines the interface for accessing job data
type JobRepository interface {
	// FetchActiveJobs retrieves all active jobs for a repository or organization
	FetchActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error)
}
