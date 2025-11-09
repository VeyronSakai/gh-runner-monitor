package entity

import (
	"testing"
	"time"
)

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

func TestJobMethods(t *testing.T) {
	t.Run("IsRunning", func(t *testing.T) {
		job := &Job{Status: "in_progress"}
		if !job.IsRunning() {
			t.Error("expected IsRunning() to be true for in_progress status")
		}

		job.Status = "queued"
		if job.IsRunning() {
			t.Error("expected IsRunning() to be false for queued status")
		}
	})

	t.Run("IsAssignedToRunner", func(t *testing.T) {
		runnerID := int64(123)
		job := &Job{RunnerID: &runnerID}

		if !job.IsAssignedToRunner(123) {
			t.Error("expected IsAssignedToRunner(123) to be true")
		}

		if job.IsAssignedToRunner(456) {
			t.Error("expected IsAssignedToRunner(456) to be false")
		}

		job.RunnerID = nil
		if job.IsAssignedToRunner(123) {
			t.Error("expected IsAssignedToRunner(123) to be false when RunnerID is nil")
		}
	})
}
