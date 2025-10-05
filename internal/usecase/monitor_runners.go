package usecase

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/service"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/value_object"
)

// MonitorRunnersUseCase handles the business logic for monitoring runners
type MonitorRunnersUseCase struct {
	runnerRepo repository.RunnerRepository
}

// NewMonitorRunnersUseCase creates a new MonitorRunnersUseCase
func NewMonitorRunnersUseCase(
	runnerRepo repository.RunnerRepository,
) *MonitorRunnersUseCase {
	return &MonitorRunnersUseCase{
		runnerRepo: runnerRepo,
	}
}

// Execute retrieves runners and jobs, and updates runner status
func (u *MonitorRunnersUseCase) Execute(ctx context.Context, owner, repo, org string) (*value_object.MonitorData, error) {
	// Fetch runners
	runners, err := u.runnerRepo.GetRunners(ctx, owner, repo, org)
	if err != nil {
		return nil, err
	}

	// Fetch active jobs
	jobs, err := u.runnerRepo.GetActiveJobs(ctx, owner, repo, org)
	if err != nil {
		return nil, err
	}

	// Update runner status based on active jobs
	service.UpdateRunnerStatus(runners, jobs)

	return &value_object.MonitorData{
		Runners: runners,
		Jobs:    jobs,
	}, nil
}
