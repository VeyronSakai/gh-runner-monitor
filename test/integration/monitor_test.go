package integration

import (
	"context"
	"testing"
	"time"
	
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/services"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
	"github.com/VeyronSakai/gh-runner-monitor/test/mocks"
)

func TestMonitorIntegration(t *testing.T) {
	ctx := context.Background()
	
	t.Run("full monitoring workflow with multiple runners", func(t *testing.T) {
		// Setup mock data
		runnerID1 := int64(1)
		runnerID2 := int64(2)
		runnerID3 := int64(3)
		startedAt := time.Now().Add(-10 * time.Minute)
		
		runners := []*entities.Runner{
			{
				ID:     runnerID1,
				Name:   "runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted", "linux", "x64"},
			},
			{
				ID:     runnerID2,
				Name:   "runner-2",
				Status: entities.StatusIdle,
				OS:     "windows",
				Labels: []string{"self-hosted", "windows", "x64"},
			},
			{
				ID:     runnerID3,
				Name:   "runner-3",
				Status: entities.StatusOffline,
				OS:     "macos",
				Labels: []string{"self-hosted", "macos", "arm64"},
			},
		}
		
		jobs := []*entities.Job{
			{
				ID:           100,
				RunID:        200,
				Name:         "Build and Test",
				Status:       entities.JobStatusInProgress,
				RunnerID:     &runnerID1,
				RunnerName:   &runners[0].Name,
				StartedAt:    &startedAt,
				WorkflowName: "CI Pipeline",
				Repository:   "owner/repo",
			},
			{
				ID:           101,
				RunID:        201,
				Name:         "Deploy",
				Status:       entities.JobStatusQueued,
				WorkflowName: "CD Pipeline",
				Repository:   "owner/repo",
			},
		}
		
		// Setup mocks
		mockRunnerRepo := mocks.NewMockRunnerRepository().
			WithListRunners(func(ctx context.Context, owner, repo string) ([]*entities.Runner, error) {
				return runners, nil
			})
		
		mockJobRepo := mocks.NewMockJobRepository().
			WithListActiveJobs(func(ctx context.Context, owner, repo string) ([]*entities.Job, error) {
				return jobs, nil
			})
		
		// Create service and use case
		monitorService := services.NewRunnerMonitorService()
		useCase := use_cases.NewMonitorRunnersUseCase(mockRunnerRepo, mockJobRepo, monitorService)
		
		// Execute
		params := ports.NewMonitorParams("owner", "repo")
		output, err := useCase.Execute(ctx, params)
		
		// Assertions
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		if len(output.Runners) != 3 {
			t.Errorf("expected 3 runners, got %d", len(output.Runners))
		}
		
		// Check runner-1 (should be Active)
		runner1 := findRunnerByName(output.Runners, "runner-1")
		if runner1 == nil {
			t.Fatal("runner-1 not found")
		}
		if runner1.Status != "Active" {
			t.Errorf("expected runner-1 to be Active, got %s", runner1.Status)
		}
		if runner1.StatusIcon != "🟠" {
			t.Errorf("expected orange icon for Active, got %s", runner1.StatusIcon)
		}
		if runner1.JobName != "Build and Test" {
			t.Errorf("expected job name 'Build and Test', got %s", runner1.JobName)
		}
		if runner1.ExecutionTime == "-" {
			t.Error("expected execution time for active job")
		}
		
		// Check runner-2 (should be Idle)
		runner2 := findRunnerByName(output.Runners, "runner-2")
		if runner2 == nil {
			t.Fatal("runner-2 not found")
		}
		if runner2.Status != "Idle" {
			t.Errorf("expected runner-2 to be Idle, got %s", runner2.Status)
		}
		if runner2.StatusIcon != "🟢" {
			t.Errorf("expected green icon for Idle, got %s", runner2.StatusIcon)
		}
		if runner2.JobName != "-" {
			t.Errorf("expected no job for idle runner, got %s", runner2.JobName)
		}
		
		// Check runner-3 (should be Offline)
		runner3 := findRunnerByName(output.Runners, "runner-3")
		if runner3 == nil {
			t.Fatal("runner-3 not found")
		}
		if runner3.Status != "Offline" {
			t.Errorf("expected runner-3 to be Offline, got %s", runner3.Status)
		}
		if runner3.StatusIcon != "⚫" {
			t.Errorf("expected gray icon for Offline, got %s", runner3.StatusIcon)
		}
	})
	
	t.Run("organization level monitoring", func(t *testing.T) {
		// Setup mock data
		runners := []*entities.Runner{
			{
				ID:     1,
				Name:   "org-runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted", "org"},
			},
		}
		
		// Setup mocks
		mockRunnerRepo := mocks.NewMockRunnerRepository().
			WithListOrgRunners(func(ctx context.Context, org string) ([]*entities.Runner, error) {
				if org != "test-org" {
					t.Errorf("expected org 'test-org', got %s", org)
				}
				return runners, nil
			})
		
		mockJobRepo := mocks.NewMockJobRepository().
			WithListOrgActiveJobs(func(ctx context.Context, org string) ([]*entities.Job, error) {
				return []*entities.Job{}, nil
			})
		
		// Create service and use case
		monitorService := services.NewRunnerMonitorService()
		useCase := use_cases.NewMonitorRunnersUseCase(mockRunnerRepo, mockJobRepo, monitorService)
		
		// Execute
		params := ports.NewOrgMonitorParams("test-org")
		output, err := useCase.Execute(ctx, params)
		
		// Assertions
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		if !params.IsOrgLevel {
			t.Error("expected IsOrgLevel to be true")
		}
		
		if len(output.Runners) != 1 {
			t.Errorf("expected 1 runner, got %d", len(output.Runners))
		}
		
		if output.Runners[0].Name != "org-runner-1" {
			t.Errorf("expected org-runner-1, got %s", output.Runners[0].Name)
		}
	})
}

func findRunnerByName(runners []*ports.RunnerOutput, name string) *ports.RunnerOutput {
	for _, r := range runners {
		if r.Name == name {
			return r
		}
	}
	return nil
}