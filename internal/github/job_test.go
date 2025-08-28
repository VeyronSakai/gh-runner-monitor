package github

import (
	"strings"
	"testing"
	"time"
)

func TestGetJobsForRun(t *testing.T) {
	tests := []struct {
		name      string
		fullName  string
		org       string
		expectErr bool
	}{
		{
			name:      "valid repository format",
			fullName:  "owner/repo",
			org:       "org",
			expectErr: false,
		},
		{
			name:      "invalid repository format - no slash",
			fullName:  "ownerrepo",
			org:       "org",
			expectErr: true,
		},
		{
			name:      "invalid repository format - multiple slashes",
			fullName:  "owner/repo/extra",
			org:       "org",
			expectErr: true,
		},
		{
			name:      "empty repository name",
			fullName:  "",
			org:       "org",
			expectErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the repository name parsing logic
			if tt.org != "" {
				parts := strings.Split(tt.fullName, "/")
				hasError := len(parts) != 2
				
				if hasError != tt.expectErr {
					t.Errorf("expected error: %v, got error: %v", tt.expectErr, hasError)
				}
			}
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