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
		validateData  func(*testing.T, []*entity.Runner, []*entity.Job)
	}{
		{
			name: "successful execution with active runners and jobs",
			runners: []*entity.Runner{
				{
					ID:        1,
					Name:      "runner-1",
					Status:    entity.StatusIdle,
					Labels:    []string{"self-hosted", "linux"},
					OS:        "Linux",
					UpdatedAt: time.Now(),
				},
				{
					ID:        2,
					Name:      "runner-2",
					Status:    entity.StatusIdle,
					Labels:    []string{"self-hosted", "macos"},
					OS:        "macOS",
					UpdatedAt: time.Now(),
				},
			},
			jobs: []*entity.Job{
				{
					ID:           1,
					RunID:        100,
					Name:         "build",
					Status:       "in_progress",
					RunnerID:     &runnerID,
					RunnerName:   &runnerName,
					StartedAt:    &startedAt,
					WorkflowName: "CI",
					Repository:   "owner/repo",
				},
			},
			wantErr: false,
			validateData: func(t *testing.T, runners []*entity.Runner, jobs []*entity.Job) {
				if len(runners) != 2 {
					t.Errorf("Expected 2 runners, got %d", len(runners))
				}
				if len(jobs) != 1 {
					t.Errorf("Expected 1 job, got %d", len(jobs))
				}
			},
		},
		{
			name:    "no runners or jobs",
			runners: []*entity.Runner{},
			jobs:    []*entity.Job{},
			wantErr: false,
			validateData: func(t *testing.T, runners []*entity.Runner, jobs []*entity.Job) {
				if len(runners) != 0 {
					t.Errorf("Expected 0 runners, got %d", len(runners))
				}
				if len(jobs) != 0 {
					t.Errorf("Expected 0 jobs, got %d", len(jobs))
				}
			},
		},
		{
			name:          "error fetching runners",
			runners:       nil,
			jobs:          nil,
			getRunnersErr: errors.New("failed to fetch runners"),
			wantErr:       true,
		},
		{
			name: "error fetching jobs",
			runners: []*entity.Runner{
				{
					ID:        1,
					Name:      "runner-1",
					Status:    entity.StatusIdle,
					Labels:    []string{"self-hosted"},
					OS:        "Linux",
					UpdatedAt: time.Now(),
				},
			},
			jobs:       nil,
			getJobsErr: errors.New("failed to fetch jobs"),
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

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if data != nil {
					t.Error("Execute() should return nil data on error")
				}
				return
			}

			if data == nil {
				t.Fatal("Execute() returned nil data")
			}

			if tt.validateData != nil {
				tt.validateData(t, data.Runners, data.Jobs)
			}
		})
	}
}

func TestRunnerMonitor_Execute_OrganizationLevel(t *testing.T) {
	repo := &test.StubRunnerRepository{
		Runners: []*entity.Runner{
			{
				ID:        1,
				Name:      "org-runner-1",
				Status:    entity.StatusIdle,
				Labels:    []string{"self-hosted"},
				OS:        "Linux",
				UpdatedAt: time.Now(),
			},
		},
		Jobs: []*entity.Job{},
	}

	useCase := NewRunnerMonitor(repo)
	ctx := context.Background()

	data, err := useCase.Execute(ctx, "", "", "my-org")

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	if data == nil {
		t.Fatal("Execute() returned nil data")
	}

	if len(data.Runners) != 1 {
		t.Errorf("Expected 1 runner, got %d", len(data.Runners))
	}
}

func TestRunnerMonitor_Execute_RunnerStatusUpdate(t *testing.T) {
	runnerID := int64(1)
	runnerName := "test-runner"
	startedAt := time.Now()

	repo := &test.StubRunnerRepository{
		Runners: []*entity.Runner{
			{
				ID:        1,
				Name:      "test-runner",
				Status:    entity.StatusIdle,
				Labels:    []string{"self-hosted"},
				OS:        "Linux",
				UpdatedAt: time.Now(),
			},
		},
		Jobs: []*entity.Job{
			{
				ID:           1,
				RunID:        100,
				Name:         "build",
				Status:       "in_progress",
				RunnerID:     &runnerID,
				RunnerName:   &runnerName,
				StartedAt:    &startedAt,
				WorkflowName: "CI",
				Repository:   "owner/repo",
			},
		},
	}

	useCase := NewRunnerMonitor(repo)
	ctx := context.Background()

	data, err := useCase.Execute(ctx, "owner", "repo", "")

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	if data == nil {
		t.Fatal("Execute() returned nil data")
	}

	// Verify that runner status was updated to Active
	if len(data.Runners) > 0 {
		runner := data.Runners[0]
		if runner.Status != entity.StatusActive {
			t.Errorf("Expected runner status to be Active, got %s", runner.Status)
		}
	}
}
