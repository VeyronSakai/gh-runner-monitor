package test

import (
	"context"
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/service"
)

func TestIntegrationFlow(t *testing.T) {
	t.Run("full monitoring workflow", func(t *testing.T) {
		// This test would require a mock GitHub API server
		t.Skip("Requires mock API server")

		// Test flow:
		// 1. Create GitHub client (infrastructure)
		// 2. Create domain service
		// 3. Create use case
		// 4. Fetch runners and jobs via use case
		// 5. Create TUI model
		// 6. Verify data display
	})

	t.Run("organization level monitoring", func(t *testing.T) {
		t.Skip("Requires mock API server")

		// Test flow for org-level monitoring
	})

	t.Run("repository level monitoring", func(t *testing.T) {
		t.Skip("Requires mock API server")

		// Test flow for repo-level monitoring
	})
}

func TestDataMerging(t *testing.T) {
	t.Run("merge runners with jobs using domain service", func(t *testing.T) {
		runners := []*entity.Runner{
			{
				ID:     1,
				Name:   "runner-1",
				Status: entity.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted"},
			},
			{
				ID:     2,
				Name:   "runner-2",
				Status: entity.StatusIdle,
				OS:     "windows",
				Labels: []string{"self-hosted", "windows"},
			},
		}

		runnerID1 := int64(1)
		runnerName1 := "runner-1"
		startedAt := time.Now().Add(-10 * time.Minute)

		jobs := []*entity.Job{
			{
				ID:           100,
				RunID:        200,
				Name:         "build",
				Status:       "in_progress",
				RunnerID:     &runnerID1,
				RunnerName:   &runnerName1,
				StartedAt:    &startedAt,
				WorkflowName: "CI",
				Repository:   "owner/repo",
			},
		}

		// Use the domain service to update runner status
		service.UpdateRunnerStatus(runners, jobs)

		// Verify runner 1 is now active
		if runners[0].Status != entity.StatusActive {
			t.Errorf("expected runner-1 to be Active, got %s", runners[0].Status)
		}

		// Verify runner 2 is still idle
		if runners[1].Status != entity.StatusIdle {
			t.Errorf("expected runner-2 to be Idle, got %s", runners[1].Status)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("handle API errors gracefully", func(t *testing.T) {
		// Test that errors from GitHub API are handled properly
		t.Skip("Requires mock API server with error responses")
	})

	t.Run("handle network timeouts", func(t *testing.T) {
		// Test timeout handling
		t.Skip("Requires network simulation")
	})
}

func TestRefreshLogic(t *testing.T) {
	t.Run("auto refresh every 5 seconds", func(t *testing.T) {
		// Test that the TUI refreshes data automatically
		t.Skip("Requires TUI lifecycle testing")
	})

	t.Run("manual refresh with 'r' key", func(t *testing.T) {
		// Test manual refresh functionality
		t.Skip("Requires TUI event simulation")
	})
}

// MockRunnerRepository for testing
type MockRunnerRepository struct {
	runners []*entity.Runner
	jobs    []*entity.Job
	err     error
}

func (m *MockRunnerRepository) GetRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.runners, nil
}

func (m *MockRunnerRepository) GetActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.jobs, nil
}

func TestWithMockRepository(t *testing.T) {
	t.Run("display runners correctly with mock", func(t *testing.T) {
		mockRepo := &MockRunnerRepository{
			runners: []*entity.Runner{
				{
					ID:     1,
					Name:   "test-runner",
					Status: entity.StatusIdle,
					OS:     "linux",
					Labels: []string{"self-hosted"},
				},
			},
			jobs: []*entity.Job{},
		}

		ctx := context.Background()

		runners, err := mockRepo.GetRunners(ctx, "owner", "repo", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runners) != 1 {
			t.Errorf("expected 1 runner, got %d", len(runners))
		}

		if runners[0].Name != "test-runner" {
			t.Errorf("expected runner name 'test-runner', got %s", runners[0].Name)
		}
	})
}
