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
	repo := &test.StubRunnerRepository{}
	useCase := NewRunnerMonitor(repo)

	if useCase == nil {
		t.Fatal("NewRunnerMonitor() returned nil")
	}

	if useCase.runnerRepo != repo {
		t.Error("NewRunnerMonitor() did not set runnerRepo correctly")
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
					Status: entity.StatusIdle,
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
					Status: entity.StatusIdle,
				},
			},
			getJobsErr: errors.New("failed to get jobs"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &test.StubRunnerRepository{
				Runners:            tt.runners,
				Jobs:               tt.jobs,
				GetRunnersError:    tt.getRunnersErr,
				GetActiveJobsError: tt.getJobsErr,
			}

			useCase := NewRunnerMonitor(repo)
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
	runnerID := int64(1)
	runnerName := "test-runner"
	startedAt := time.Now()

	repo := &test.StubRunnerRepository{
		Runners: []*entity.Runner{
			{
				ID:     runnerID,
				Name:   runnerName,
				Status: entity.StatusIdle,
				Labels: []string{"self-hosted"},
				OS:     "linux",
			},
		},
		Jobs: []*entity.Job{
			{
				ID:         1,
				RunID:      100,
				Name:       "job-1",
				Status:     "in_progress",
				RunnerID:   &runnerID,
				RunnerName: &runnerName,
				StartedAt:  &startedAt,
			},
		},
	}

	useCase := NewRunnerMonitor(repo)
	ctx := context.Background()

	data, err := useCase.Execute(ctx, "owner", "repo", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Runner should be marked as Active because it has an active job
	if data.Runners[0].Status != entity.StatusActive {
		t.Errorf("Expected runner status Active, got %s", data.Runners[0].Status)
	}
}
