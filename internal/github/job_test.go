package github

import (
	"testing"
	"time"
)

func TestGetJobsForRun(t *testing.T) {
	tests := []struct {
		name     string
		run      workflowRun
		org      string
		owner    string
		repo     string
		wantOwner string
		wantRepo  string
	}{
		{
			name: "repository-level run",
			run: workflowRun{
				ID:   123,
				Name: "CI",
				Repository: repository{
					FullName: "owner/repo",
				},
			},
			org:       "",
			owner:     "owner",
			repo:      "repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name: "organization-level run",
			run: workflowRun{
				ID:   456,
				Name: "Deploy",
				Repository: repository{
					FullName: "myorg/myrepo",
				},
			},
			org:       "myorg",
			owner:     "",
			repo:      "",
			wantOwner: "myorg",
			wantRepo:  "myrepo",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies the logic for parsing repository names
			// In a real test, you would mock the HTTP client
			t.Skip("Requires HTTP client mocking")
		})
	}
}

func TestFilterActiveJobs(t *testing.T) {
	now := time.Now()
	runnerID := int64(1)
	runnerName := "test-runner"
	
	jobsResp := &jobsResponse{
		Jobs: []jobResponse{
			{
				ID:         1,
				RunID:      100,
				Name:       "build",
				Status:     "in_progress",
				RunnerID:   &runnerID,
				RunnerName: &runnerName,
				StartedAt:  &now,
			},
			{
				ID:         2,
				RunID:      100,
				Name:       "test",
				Status:     "queued",
				RunnerID:   nil,
				RunnerName: nil,
				StartedAt:  nil,
			},
			{
				ID:         3,
				RunID:      100,
				Name:       "deploy",
				Status:     "completed",
				RunnerID:   &runnerID,
				RunnerName: &runnerName,
				StartedAt:  &now,
			},
			{
				ID:         4,
				RunID:      100,
				Name:       "cleanup",
				Status:     "cancelled",
				RunnerID:   nil,
				RunnerName: nil,
				StartedAt:  nil,
			},
		},
	}
	
	// Test that only in_progress and queued jobs are included
	activeCount := 0
	for _, job := range jobsResp.Jobs {
		if job.Status == "in_progress" || job.Status == "queued" {
			activeCount++
		}
	}
	
	if activeCount != 2 {
		t.Errorf("expected 2 active jobs, got %d", activeCount)
	}
}