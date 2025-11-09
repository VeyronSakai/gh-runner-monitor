package repository

import (
	"context"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// RunnerRepository defines the interface for accessing runner data
type RunnerRepository interface {
	// FetchRunners all runners for a repository or organization
	FetchRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error)

	// FetchActiveJobs retrieves all active jobs for a repository or organization
	FetchActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error)

	// GetCurrentTime returns the current time (for mocking in tests/debug mode)
	GetCurrentTime() time.Time
}
