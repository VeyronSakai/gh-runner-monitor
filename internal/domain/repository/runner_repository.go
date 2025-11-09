package repository

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// RunnerRepository defines the interface for accessing runner data
type RunnerRepository interface {
	// FetchRunners all runners for a repository or organization
	FetchRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error)
}
