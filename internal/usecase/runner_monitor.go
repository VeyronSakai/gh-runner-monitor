package usecase

import (
	"context"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/service"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/value_object"
)

// RunnerMonitor handles the business logic for monitoring runners
type RunnerMonitor struct {
	runnerRepo repository.RunnerRepository
}

// NewRunnerMonitor creates a new RunnerMonitor
func NewRunnerMonitor(runnerRepo repository.RunnerRepository) *RunnerMonitor {
	return &RunnerMonitor{
		runnerRepo: runnerRepo,
	}
}

// Execute retrieves runners and jobs, and updates runner status
func (u *RunnerMonitor) Execute(ctx context.Context, owner, repo, org string) (*value_object.MonitorData, error) {
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

	// Get current time (use mocked time if available, otherwise use actual time)
	currentTime := time.Now()
	if timeProvider, ok := u.runnerRepo.(interface{ GetCurrentTime() time.Time }); ok {
		currentTime = timeProvider.GetCurrentTime()
	}

	return &value_object.MonitorData{
		CurrentTime: currentTime,
		Runners:     runners,
		Jobs:        jobs,
	}, nil
}
