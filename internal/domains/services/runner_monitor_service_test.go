package services

import (
	"testing"
	"time"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

func TestRunnerMonitorService_MergeRunnersAndJobs(t *testing.T) {
	service := NewRunnerMonitorService()
	
	t.Run("merge runners and active jobs", func(t *testing.T) {
		runnerID := int64(1)
		startedAt := time.Now().Add(-5 * time.Minute)
		
		runners := []*entities.Runner{
			{
				ID:     runnerID,
				Name:   "runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted"},
			},
			{
				ID:     2,
				Name:   "runner-2",
				Status: entities.StatusOffline,
				OS:     "windows",
				Labels: []string{"self-hosted"},
			},
		}
		
		jobs := []*entities.Job{
			{
				ID:         100,
				RunID:      200,
				Name:       "Build",
				Status:     entities.JobStatusInProgress,
				RunnerID:   &runnerID,
				StartedAt:  &startedAt,
			},
		}
		
		result := service.MergeRunnersAndJobs(runners, jobs)
		
		if len(result) != 2 {
			t.Errorf("expected 2 results, got %d", len(result))
		}
		
		activeRunner := result[0]
		if activeRunner.Runner.Status != entities.StatusActive {
			t.Errorf("expected runner-1 to be Active, got %s", activeRunner.Runner.Status)
		}
		if activeRunner.CurrentJob == nil {
			t.Error("expected runner-1 to have a job")
		}
		if activeRunner.ExecutionTime == 0 {
			t.Error("expected runner-1 to have execution time")
		}
		
		offlineRunner := result[1]
		if offlineRunner.Runner.Status != entities.StatusOffline {
			t.Errorf("expected runner-2 to be Offline, got %s", offlineRunner.Runner.Status)
		}
		if offlineRunner.CurrentJob != nil {
			t.Error("expected runner-2 to have no job")
		}
	})
	
	t.Run("all runners idle when no jobs", func(t *testing.T) {
		runners := []*entities.Runner{
			{
				ID:     1,
				Name:   "runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
			},
			{
				ID:     2,
				Name:   "runner-2",
				Status: entities.StatusIdle,
				OS:     "windows",
			},
		}
		
		jobs := []*entities.Job{}
		
		result := service.MergeRunnersAndJobs(runners, jobs)
		
		if len(result) != 2 {
			t.Errorf("expected 2 results, got %d", len(result))
		}
		
		for _, rwj := range result {
			if rwj.Runner.Status != entities.StatusIdle {
				t.Errorf("expected all runners to be Idle, got %s for %s", 
					rwj.Runner.Status, rwj.Runner.Name)
			}
		}
	})
}

func TestRunnerMonitorService_DetermineRunnerStatus(t *testing.T) {
	service := NewRunnerMonitorService()
	
	t.Run("offline runner stays offline", func(t *testing.T) {
		runner := &entities.Runner{
			ID:     1,
			Name:   "runner-1",
			Status: entities.StatusOffline,
		}
		
		activeJob := &entities.Job{
			Status: entities.JobStatusInProgress,
		}
		
		status := service.DetermineRunnerStatus(runner, activeJob)
		
		if status != entities.StatusOffline {
			t.Errorf("expected Offline, got %s", status)
		}
	})
	
	t.Run("runner with active job becomes active", func(t *testing.T) {
		runner := &entities.Runner{
			ID:     1,
			Name:   "runner-1",
			Status: entities.StatusIdle,
		}
		
		activeJob := &entities.Job{
			Status: entities.JobStatusInProgress,
		}
		
		status := service.DetermineRunnerStatus(runner, activeJob)
		
		if status != entities.StatusActive {
			t.Errorf("expected Active, got %s", status)
		}
	})
	
	t.Run("runner without job is idle", func(t *testing.T) {
		runner := &entities.Runner{
			ID:     1,
			Name:   "runner-1",
			Status: entities.StatusIdle,
		}
		
		status := service.DetermineRunnerStatus(runner, nil)
		
		if status != entities.StatusIdle {
			t.Errorf("expected Idle, got %s", status)
		}
	})
}