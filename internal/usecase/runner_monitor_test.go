package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/test"
)

func TestNewRunnerMonitor(t *testing.T) {
	runnerRepo := &test.StubRunnerRepository{}
	jobRepo := &test.StubJobRepository{}
	timeProvider := &test.StubTimeProvider{}
	useCase := NewRunnerMonitor(runnerRepo, jobRepo, timeProvider)

	if useCase == nil {
		t.Fatal("NewRunnerMonitor() returned nil")
	}

	if useCase.runnerRepo != runnerRepo {
		t.Error("NewRunnerMonitor() did not set runnerRepo correctly")
	}
	if useCase.jobRepo != jobRepo {
		t.Error("NewRunnerMonitor() did not set jobRepo correctly")
	}
	if useCase.timeProvider != timeProvider {
		t.Error("NewRunnerMonitor() did not set timeProvider correctly")
	}
}

func TestRunnerMonitor_Execute(t *testing.T) {
	runnerID := int64(1)
	runnerName := "test-runner"
	startedAt := time.Now()

	tests := []struct {
		name          string
		runners       []*entity.Runner
		jobs          []*entity.Job
		getRunnersErr error
		getJobsErr    error
		wantErr       bool
		validateData  func(*testing.T, []*entity.Runner, []*entity.Job, time.Time)
	}{
		{
			name: "successful execution with active runners and jobs",
			runners: []*entity.Runner{
				{
					ID:     runnerID,
					Name:   runnerName,
					Status: entity.StatusActive,
					Labels: []string{"self-hosted", "linux"},
					OS:     "linux",
				},
			},
			jobs: []*entity.Job{
				{
					ID:         1,
					RunID:      100,
					Name:       "test-job",
					Status:     "in_progress",
					RunnerID:   &runnerID,
					RunnerName: &runnerName,
					StartedAt:  &startedAt,
				},
			},
			wantErr: false,
			validateData: func(t *testing.T, runners []*entity.Runner, jobs []*entity.Job, currentTime time.Time) {
				if len(runners) != 1 {
					t.Errorf("Expected 1 runner, got %d", len(runners))
				}
				if len(jobs) != 1 {
					t.Errorf("Expected 1 job, got %d", len(jobs))
				}
				// Verify that runner status was updated to Active
				if runners[0].Status != entity.StatusActive {
					t.Errorf("Expected runner status to be Active, got %s", runners[0].Status)
				}
				// Verify current time is set
				if currentTime.IsZero() {
					t.Error("Expected current time to be set")
				}
			},
		},
		{
			name:    "no runners and no jobs",
			runners: []*entity.Runner{},
			jobs:    []*entity.Job{},
			wantErr: false,
			validateData: func(t *testing.T, runners []*entity.Runner, jobs []*entity.Job, currentTime time.Time) {
				if len(runners) != 0 {
					t.Errorf("Expected 0 runners, got %d", len(runners))
				}
				if len(jobs) != 0 {
					t.Errorf("Expected 0 jobs, got %d", len(jobs))
				}
			},
		},
		{
			name:          "FetchRunners returns error",
			getRunnersErr: errors.New("failed to get runners"),
			wantErr:       true,
		},
		{
			name: "FetchActiveJobs returns error",
			runners: []*entity.Runner{
				{
					ID:     runnerID,
					Name:   runnerName,
					Status: entity.StatusActive,
				},
			},
			getJobsErr: errors.New("failed to get jobs"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runnerRepo := &test.StubRunnerRepository{
				Runners:         tt.runners,
				GetRunnersError: tt.getRunnersErr,
			}
			jobRepo := &test.StubJobRepository{
				Jobs:               tt.jobs,
				GetActiveJobsError: tt.getJobsErr,
			}
			timeProvider := &test.StubTimeProvider{}

			useCase := NewRunnerMonitor(runnerRepo, jobRepo, timeProvider)
			ctx := context.Background()

			data, err := useCase.Execute(ctx, "owner", "repo", "")

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if data == nil {
				t.Fatal("Expected data but got nil")
			}

			if tt.validateData != nil {
				tt.validateData(t, data.Runners, data.Jobs, data.CurrentTime)
			}
		})
	}
}

func TestRunnerMonitor_Execute_UpdatesRunnerStatus(t *testing.T) {
	busyRunnerID := int64(1)
	idleRunnerID := int64(2)
	idleRunnerName := "idle-runner"
	startedAt := time.Now()

	runnerRepo := &test.StubRunnerRepository{
		Runners: []*entity.Runner{
			{ID: busyRunnerID, Name: "busy-runner", Status: entity.StatusActive},
			// Reported as idle by the API, but a job has just been assigned to it
			{ID: idleRunnerID, Name: idleRunnerName, Status: entity.StatusIdle},
		},
	}
	jobRepo := &test.StubJobRepository{
		Jobs: []*entity.Job{
			{
				ID:         1,
				RunID:      100,
				Name:       "job-1",
				Status:     "in_progress",
				RunnerID:   &idleRunnerID,
				RunnerName: &idleRunnerName,
				StartedAt:  &startedAt,
			},
		},
	}
	timeProvider := &test.StubTimeProvider{}

	useCase := NewRunnerMonitor(runnerRepo, jobRepo, timeProvider)
	ctx := context.Background()

	data, err := useCase.Execute(ctx, "owner", "repo", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Runner should be marked as Active because it has an active job
	if data.Runners[1].Status != entity.StatusActive {
		t.Errorf("Expected runner status Active, got %s", data.Runners[1].Status)
	}
}

func TestRunnerMonitor_Execute_SkipsJobFetchWhenNoRunnerIsBusy(t *testing.T) {
	runnerRepo := &test.StubRunnerRepository{
		Runners: []*entity.Runner{
			{ID: 1, Name: "idle-runner", Status: entity.StatusIdle},
			{ID: 2, Name: "offline-runner", Status: entity.StatusOffline},
		},
	}
	jobRepo := &test.StubJobRepository{
		GetActiveJobsError: errors.New("should not be called"),
	}
	timeProvider := &test.StubTimeProvider{}

	useCase := NewRunnerMonitor(runnerRepo, jobRepo, timeProvider)

	data, err := useCase.Execute(context.Background(), "", "", "org")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if jobRepo.CallCount != 0 {
		t.Errorf("Expected FetchActiveJobs not to be called, got %d calls", jobRepo.CallCount)
	}
	if len(data.Jobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(data.Jobs))
	}
}
