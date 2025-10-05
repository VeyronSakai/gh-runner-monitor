package service

import (
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// RunnerService provides domain logic for runners
type RunnerService struct{}

// NewRunnerService creates a new RunnerService
func NewRunnerService() *RunnerService {
	return &RunnerService{}
}

// UpdateRunnerStatus updates the runner status based on active jobs
// If a runner has an active job, its status is set to Active
func (s *RunnerService) UpdateRunnerStatus(runners []*entity.Runner, jobs []*entity.Job) {
	for _, runner := range runners {
		// Reset to idle if the runner is online but not assigned to any job
		if runner.IsOnline() && !runner.IsActive() {
			runner.Status = entity.StatusIdle
		}

		// Check if this runner has an active job
		for _, job := range jobs {
			if job.IsAssignedToRunner(runner.ID) && job.IsRunning() {
				runner.Status = entity.StatusActive
				break
			}
		}
	}
}

// FindJobForRunner finds the active job assigned to a specific runner
func (s *RunnerService) FindJobForRunner(runnerID int64, jobs []*entity.Job) *entity.Job {
	for _, job := range jobs {
		if job.IsAssignedToRunner(runnerID) {
			return job
		}
	}
	return nil
}
