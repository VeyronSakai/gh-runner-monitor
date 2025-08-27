package use_cases

import (
	"context"
	"fmt"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/repositories"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/services"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
	"time"
)

type MonitorRunnersUseCase struct {
	runnerRepo   repositories.RunnerRepository
	jobRepo      repositories.JobRepository
	monitorService *services.RunnerMonitorService
}

func NewMonitorRunnersUseCase(
	runnerRepo repositories.RunnerRepository,
	jobRepo repositories.JobRepository,
	monitorService *services.RunnerMonitorService,
) *MonitorRunnersUseCase {
	return &MonitorRunnersUseCase{
		runnerRepo:   runnerRepo,
		jobRepo:      jobRepo,
		monitorService: monitorService,
	}
}

func (uc *MonitorRunnersUseCase) Execute(ctx context.Context, params ports.MonitorParams) (*ports.MonitorOutput, error) {
	var runners []*entities.Runner
	var jobs []*entities.Job
	var err error
	
	if params.IsOrgLevel {
		runners, err = uc.runnerRepo.ListOrgRunners(ctx, params.Org)
		if err != nil {
			return nil, fmt.Errorf("failed to list org runners: %w", err)
		}
		
		jobs, err = uc.jobRepo.ListOrgActiveJobs(ctx, params.Org)
		if err != nil {
			return nil, fmt.Errorf("failed to list org active jobs: %w", err)
		}
	} else {
		runners, err = uc.runnerRepo.ListRunners(ctx, params.Owner, params.Repo)
		if err != nil {
			return nil, fmt.Errorf("failed to list runners: %w", err)
		}
		
		jobs, err = uc.jobRepo.ListActiveJobs(ctx, params.Owner, params.Repo)
		if err != nil {
			return nil, fmt.Errorf("failed to list active jobs: %w", err)
		}
	}
	
	runnersWithJobs := uc.monitorService.MergeRunnersAndJobs(runners, jobs)
	
	output := &ports.MonitorOutput{
		Runners:   make([]*ports.RunnerOutput, 0, len(runnersWithJobs)),
		UpdatedAt: time.Now(),
	}
	
	for _, rwj := range runnersWithJobs {
		status := uc.monitorService.DetermineRunnerStatus(rwj.Runner, rwj.CurrentJob)
		runnerOutput := ports.NewRunnerOutput(rwj.Runner.ID, rwj.Runner.Name, status.String())
		runnerOutput.Labels = rwj.Runner.Labels
		runnerOutput.OS = rwj.Runner.OS
		
		if rwj.CurrentJob != nil {
			runnerOutput.SetJob(rwj.CurrentJob.Name, rwj.ExecutionTime)
		} else {
			runnerOutput.JobName = "-"
			runnerOutput.ExecutionTime = "-"
		}
		
		output.Runners = append(output.Runners, runnerOutput)
	}
	
	return output, nil
}