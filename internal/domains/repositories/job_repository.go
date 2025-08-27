package repositories

import (
	"context"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

type JobRepository interface {
	ListActiveJobs(ctx context.Context, owner, repo string) ([]*entities.Job, error)
	ListOrgActiveJobs(ctx context.Context, org string) ([]*entities.Job, error)
	GetJob(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error)
}