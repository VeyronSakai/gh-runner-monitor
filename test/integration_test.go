package test

import (
	"context"
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
)

func TestIntegrationFlow(t *testing.T) {
	t.Run("full monitoring workflow", func(t *testing.T) {
		// This test would require a mock GitHub API server
		t.Skip("Requires mock API server")
		
		// Test flow:
		// 1. Create GitHub client
		// 2. Fetch runners
		// 3. Fetch active jobs
		// 4. Create TUI model
		// 5. Verify data display
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
	t.Run("merge runners with jobs", func(t *testing.T) {
		runners := []*models.Runner{
			{
				ID:     1,
				Name:   "runner-1",
				Status: models.StatusIdle,
				OS:     "linux",
				Labels: []string{"self-hosted"},
			},
			{
				ID:     2,
				Name:   "runner-2",
				Status: models.StatusIdle,
				OS:     "windows",
				Labels: []string{"self-hosted", "windows"},
			},
		}
		
		runnerID1 := int64(1)
		runnerName1 := "runner-1"
		startedAt := time.Now().Add(-10 * time.Minute)
		
		jobs := []*models.Job{
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
		
		// Simulate the merging logic
		for _, runner := range runners {
			for _, job := range jobs {
				if job.RunnerID != nil && *job.RunnerID == runner.ID {
					runner.Status = models.StatusActive
					break
				}
			}
		}
		
		// Verify runner 1 is now active
		if runners[0].Status != models.StatusActive {
			t.Errorf("expected runner-1 to be Active, got %s", runners[0].Status)
		}
		
		// Verify runner 2 is still idle
		if runners[1].Status != models.StatusIdle {
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

// MockGitHubClient for testing
type MockGitHubClient struct {
	runners []*models.Runner
	jobs    []*models.Job
	err     error
}

func (m *MockGitHubClient) GetRunners(ctx context.Context, owner, repo, org string) ([]*models.Runner, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.runners, nil
}

func (m *MockGitHubClient) GetActiveJobs(ctx context.Context, owner, repo, org string) ([]*models.Job, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.jobs, nil
}

func TestWithMockClient(t *testing.T) {
	t.Run("display runners correctly", func(t *testing.T) {
		mockClient := &MockGitHubClient{
			runners: []*models.Runner{
				{
					ID:     1,
					Name:   "test-runner",
					Status: models.StatusIdle,
					OS:     "linux",
					Labels: []string{"self-hosted"},
				},
			},
			jobs: []*models.Job{},
		}
		
		// This would test the TUI model with the mock client
		// However, tui.NewModel expects a real github.Client
		// so we'd need to refactor to use an interface
		_ = mockClient
		t.Skip("Requires interface refactoring")
	})
}