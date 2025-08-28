package use_cases

import (
	"context"
	"testing"
	"time"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/services"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
)

type mockRunnerRepository struct {
	runners []*entities.Runner
	err     error
}

func (m *mockRunnerRepository) ListRunners(ctx context.Context, owner, repo string) ([]*entities.Runner, error) {
	return m.runners, m.err
}

func (m *mockRunnerRepository) ListOrgRunners(ctx context.Context, org string) ([]*entities.Runner, error) {
	return m.runners, m.err
}

func (m *mockRunnerRepository) GetRunner(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error) {
	for _, r := range m.runners {
		if r.ID == runnerID {
			return r, nil
		}
	}
	return nil, nil
}

type mockJobRepository struct {
	jobs []*entities.Job
	err  error
}

func (m *mockJobRepository) ListActiveJobs(ctx context.Context, owner, repo string) ([]*entities.Job, error) {
	return m.jobs, m.err
}

func (m *mockJobRepository) ListOrgActiveJobs(ctx context.Context, org string) ([]*entities.Job, error) {
	return m.jobs, m.err
}

func (m *mockJobRepository) GetJob(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error) {
	for _, j := range m.jobs {
		if j.ID == jobID {
			return j, nil
		}
	}
	return nil, nil
}

func TestMonitorRunnersUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	
	t.Run("successful execution with idle runner", func(t *testing.T) {
		runners := []*entities.Runner{
			{
				ID:     1,
				Name:   "runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted", "linux"},
			},
		}
		
		mockRunnerRepo := &mockRunnerRepository{runners: runners}
		mockJobRepo := &mockJobRepository{jobs: []*entities.Job{}}
		monitorService := services.NewRunnerMonitorService()
		
		useCase := NewMonitorRunnersUseCase(mockRunnerRepo, mockJobRepo, monitorService)
		
		params := ports.NewMonitorParams("owner", "repo")
		output, err := useCase.Execute(ctx, params)
		
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		if len(output.Runners) != 1 {
			t.Fatalf("expected 1 runner, got %d", len(output.Runners))
		}
		
		runner := output.Runners[0]
		if runner.Status != "Idle" {
			t.Errorf("expected status Idle, got %s", runner.Status)
		}
		if runner.StatusIcon != "🟢" {
			t.Errorf("expected green icon, got %s", runner.StatusIcon)
		}
	})
	
	t.Run("successful execution with active runner", func(t *testing.T) {
		runnerID := int64(1)
		startedAt := time.Now().Add(-5 * time.Minute)
		
		runners := []*entities.Runner{
			{
				ID:     runnerID,
				Name:   "runner-1",
				Status: entities.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted", "linux"},
			},
		}
		
		jobs := []*entities.Job{
			{
				ID:         100,
				RunID:      200,
				Name:       "Build and Test",
				Status:     entities.JobStatusInProgress,
				RunnerID:   &runnerID,
				StartedAt:  &startedAt,
			},
		}
		
		mockRunnerRepo := &mockRunnerRepository{runners: runners}
		mockJobRepo := &mockJobRepository{jobs: jobs}
		monitorService := services.NewRunnerMonitorService()
		
		useCase := NewMonitorRunnersUseCase(mockRunnerRepo, mockJobRepo, monitorService)
		
		params := ports.NewMonitorParams("owner", "repo")
		output, err := useCase.Execute(ctx, params)
		
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		if len(output.Runners) != 1 {
			t.Fatalf("expected 1 runner, got %d", len(output.Runners))
		}
		
		runner := output.Runners[0]
		if runner.Status != "Active" {
			t.Errorf("expected status Active, got %s", runner.Status)
		}
		if runner.StatusIcon != "🟠" {
			t.Errorf("expected orange icon, got %s", runner.StatusIcon)
		}
		if runner.JobName != "Build and Test" {
			t.Errorf("expected job name 'Build and Test', got %s", runner.JobName)
		}
	})
	
	t.Run("successful execution with offline runner", func(t *testing.T) {
		runners := []*entities.Runner{
			{
				ID:     1,
				Name:   "runner-1",
				Status: entities.StatusOffline,
				OS:     "linux",
				Labels: []string{"self-hosted", "linux"},
			},
		}
		
		mockRunnerRepo := &mockRunnerRepository{runners: runners}
		mockJobRepo := &mockJobRepository{jobs: []*entities.Job{}}
		monitorService := services.NewRunnerMonitorService()
		
		useCase := NewMonitorRunnersUseCase(mockRunnerRepo, mockJobRepo, monitorService)
		
		params := ports.NewMonitorParams("owner", "repo")
		output, err := useCase.Execute(ctx, params)
		
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		runner := output.Runners[0]
		if runner.Status != "Offline" {
			t.Errorf("expected status Offline, got %s", runner.Status)
		}
		if runner.StatusIcon != "⚫" {
			t.Errorf("expected gray icon, got %s", runner.StatusIcon)
		}
	})
}