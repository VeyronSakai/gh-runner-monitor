package services

import (
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
	"time"
)

type RunnerMonitorService struct{}

func NewRunnerMonitorService() *RunnerMonitorService {
	return &RunnerMonitorService{}
}

type RunnerWithJob struct {
	Runner       *entities.Runner
	CurrentJob   *entities.Job
	ExecutionTime time.Duration
}

func (s *RunnerMonitorService) MergeRunnersAndJobs(runners []*entities.Runner, jobs []*entities.Job) []*RunnerWithJob {
	runnerMap := make(map[int64]*entities.Runner)
	for _, runner := range runners {
		runnerMap[runner.ID] = runner
	}
	
	runnerWithJobs := make([]*RunnerWithJob, 0, len(runners))
	processedRunners := make(map[int64]bool)
	
	for _, job := range jobs {
		if job.IsRunning() && job.GetRunnerID() > 0 {
			if runner, ok := runnerMap[job.GetRunnerID()]; ok {
				runner.UpdateStatus(entities.StatusActive)
				runnerWithJobs = append(runnerWithJobs, &RunnerWithJob{
					Runner:        runner,
					CurrentJob:    job,
					ExecutionTime: job.GetExecutionTime(),
				})
				processedRunners[runner.ID] = true
			}
		}
	}
	
	for _, runner := range runners {
		if !processedRunners[runner.ID] {
			status := entities.StatusIdle
			if runner.Status == entities.StatusOffline {
				status = entities.StatusOffline
			}
			runner.UpdateStatus(status)
			runnerWithJobs = append(runnerWithJobs, &RunnerWithJob{
				Runner:        runner,
				CurrentJob:    nil,
				ExecutionTime: 0,
			})
		}
	}
	
	return runnerWithJobs
}

func (s *RunnerMonitorService) DetermineRunnerStatus(runner *entities.Runner, activeJob *entities.Job) entities.RunnerStatus {
	if runner.Status == entities.StatusOffline {
		return entities.StatusOffline
	}
	if activeJob != nil && activeJob.IsRunning() {
		return entities.StatusActive
	}
	return entities.StatusIdle
}