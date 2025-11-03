package repository

import (
	"context"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// RunnerRepository defines the interface for accessing runner data
type RunnerRepository interface {
	// GetRunners retrieves all runners for a repository or organization
	GetRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error)

	// GetActiveJobs retrieves all active jobs for a repository or organization
	GetActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error)

	// GetCurrentTime returns the current time (for mocking in tests/debug mode)
	GetCurrentTime() time.Time
}
