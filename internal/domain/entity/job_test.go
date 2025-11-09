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
	t.Run("IsAssignedToRunner", func(t *testing.T) {
		runnerID := int64(123)
		job := &Job{
			RunnerID: &runnerID,
			Status:   "in_progress",
		}

		if !job.IsAssignedToRunner(123) {
			t.Error("expected IsAssignedToRunner(123) to be true for running job")
		}

		if job.IsAssignedToRunner(456) {
			t.Error("expected IsAssignedToRunner(456) to be false for different runner")
		}

		job.RunnerID = nil
		if job.IsAssignedToRunner(123) {
			t.Error("expected IsAssignedToRunner(123) to be false when RunnerID is nil")
		}

		job.RunnerID = &runnerID
		job.Status = "queued"
		if job.IsAssignedToRunner(123) {
			t.Error("expected IsAssignedToRunner(123) to be false for queued job")
		}
	})

	t.Run("GetExecutionDuration", func(t *testing.T) {
		startedAt := time.Now().Add(-5 * time.Minute)
		job := &Job{StartedAt: &startedAt}

		duration := job.GetExecutionDuration()
		if duration < 4*time.Minute || duration > 6*time.Minute {
			t.Errorf("expected duration around 5 minutes, got %v", duration)
		}

		job.StartedAt = nil
		duration = job.GetExecutionDuration()
		if duration != 0 {
			t.Errorf("expected duration to be 0 when StartedAt is nil, got %v", duration)
		}
	})
}
