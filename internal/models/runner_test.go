package models

import (
	"testing"
	"time"
)

func TestRunnerStatus(t *testing.T) {
	tests := []struct {
		name   string
		status RunnerStatus
		valid  bool
	}{
		{
			name:   "idle status is valid",
			status: StatusIdle,
			valid:  true,
		},
		{
			name:   "active status is valid",
			status: StatusActive,
			valid:  true,
		},
		{
			name:   "offline status is valid",
			status: StatusOffline,
			valid:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that the status can be used as a string
			str := string(tt.status)
			if str == "" && tt.valid {
				t.Errorf("expected non-empty string for valid status")
			}
		})
	}
}

func TestRunner(t *testing.T) {
	t.Run("create new runner", func(t *testing.T) {
		runner := &Runner{
			ID:        1,
			Name:      "test-runner",
			Status:    StatusIdle,
			Labels:    []string{"self-hosted", "linux"},
			OS:        "linux",
			UpdatedAt: time.Now(),
		}
		
		if runner.ID != 1 {
			t.Errorf("expected ID 1, got %d", runner.ID)
		}
		
		if runner.Name != "test-runner" {
			t.Errorf("expected name 'test-runner', got %s", runner.Name)
		}
		
		if runner.Status != StatusIdle {
			t.Errorf("expected status Idle, got %s", runner.Status)
		}
		
		if len(runner.Labels) != 2 {
			t.Errorf("expected 2 labels, got %d", len(runner.Labels))
		}
	})
}

func TestJob(t *testing.T) {
	t.Run("create new job", func(t *testing.T) {
		now := time.Now()
		runnerID := int64(1)
		runnerName := "test-runner"
		
		job := &Job{
			ID:           100,
			RunID:        200,
			Name:         "build",
			Status:       "in_progress",
			RunnerID:     &runnerID,
			RunnerName:   &runnerName,
			StartedAt:    &now,
			WorkflowName: "CI",
			Repository:   "owner/repo",
		}
		
		if job.ID != 100 {
			t.Errorf("expected ID 100, got %d", job.ID)
		}
		
		if job.RunnerID == nil || *job.RunnerID != 1 {
			t.Errorf("expected RunnerID 1, got %v", job.RunnerID)
		}
		
		if job.StartedAt == nil {
			t.Error("expected StartedAt to be set")
		}
	})
	
	t.Run("job without runner assignment", func(t *testing.T) {
		job := &Job{
			ID:           101,
			RunID:        201,
			Name:         "test",
			Status:       "queued",
			RunnerID:     nil,
			RunnerName:   nil,
			StartedAt:    nil,
			WorkflowName: "CI",
			Repository:   "owner/repo",
		}
		
		if job.RunnerID != nil {
			t.Errorf("expected RunnerID to be nil, got %v", job.RunnerID)
		}
		
		if job.StartedAt != nil {
			t.Errorf("expected StartedAt to be nil, got %v", job.StartedAt)
		}
	})
}