package repositories

import (
	"context"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

type RunnerRepository interface {
	ListRunners(ctx context.Context, owner, repo string) ([]*entities.Runner, error)
	ListOrgRunners(ctx context.Context, org string) ([]*entities.Runner, error)
	GetRunner(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error)
}