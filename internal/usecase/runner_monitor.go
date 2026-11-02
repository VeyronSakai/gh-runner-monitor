package usecase

import (
	"context"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/service"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/value_object"
)

// RunnerMonitor handles the business logic for monitoring runners
type RunnerMonitor struct {
	runnerRepo   repository.RunnerRepository
	jobRepo      repository.JobRepository
	timeProvider repository.TimeProvider
}

// NewRunnerMonitor creates a new RunnerMonitor
func NewRunnerMonitor(
	runnerRepo repository.RunnerRepository,
	jobRepo repository.JobRepository,
	timeProvider repository.TimeProvider,
) *RunnerMonitor {
	return &RunnerMonitor{
		runnerRepo:   runnerRepo,
		jobRepo:      jobRepo,
		timeProvider: timeProvider,
	}
}

// Execute retrieves runners and jobs, and updates runner status
func (u *RunnerMonitor) Execute(ctx context.Context, owner, repo, org string) (*value_object.MonitorData, error) {
	// Fetch runners
	runners, err := u.runnerRepo.FetchRunners(ctx, owner, repo, org)
	if err != nil {
		return nil, err
	}

	// Jobs are only used to show what busy runners are doing, so skip the
	// (potentially expensive) job lookup when no runner is busy.
	var jobs []*entity.Job
	if hasActiveRunner(runners) {
		jobs, err = u.jobRepo.FetchActiveJobs(ctx, owner, repo, org)
		if err != nil {
			return nil, err
		}
	}

	// Update runner status based on active jobs
	service.UpdateRunnerStatus(runners, jobs)

	return &value_object.MonitorData{
		CurrentTime: u.timeProvider.GetCurrentTime(),
		Runners:     runners,
		Jobs:        jobs,
	}, nil
}

// hasActiveRunner reports whether any runner is currently busy
func hasActiveRunner(runners []*entity.Runner) bool {
	for _, r := range runners {
		if r.IsActive() {
			return true
		}
	}
	return false
}
