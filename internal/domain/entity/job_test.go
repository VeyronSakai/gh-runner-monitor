package entity

import (
	"testing"
	"time"
)

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

	t.Run("IsQueued", func(t *testing.T) {
		job := &Job{Status: "queued"}
		if !job.IsQueued() {
			t.Error("expected IsQueued() to be true for queued status")
		}

		job.Status = "in_progress"
		if job.IsQueued() {
			t.Error("expected IsQueued() to be false for in_progress status")
		}
	})

	t.Run("IsAssignedToRunner", func(t *testing.T) {
		runnerName := "test-runner"
		job := &Job{RunnerName: &runnerName}

		if !job.IsAssignedToRunner("test-runner") {
			t.Error("expected IsAssignedToRunner('test-runner') to be true")
		}

		if job.IsAssignedToRunner("other-runner") {
			t.Error("expected IsAssignedToRunner('other-runner') to be false")
		}

		job.RunnerName = nil
		if job.IsAssignedToRunner("test-runner") {
			t.Error("expected IsAssignedToRunner('test-runner') to be false when RunnerName is nil")
		}
	})

	t.Run("GetExecutionDurationAt", func(t *testing.T) {
		t.Run("returns duration when StartedAt is set", func(t *testing.T) {
			startTime := time.Date(2025, 11, 8, 10, 0, 0, 0, time.UTC)
			job := &Job{StartedAt: &startTime}

			currentTime := time.Date(2025, 11, 8, 10, 5, 30, 0, time.UTC)
			duration := job.GetExecutionDurationAt(currentTime)

			expected := 5*time.Minute + 30*time.Second
			if duration != expected {
				t.Errorf("expected duration %v, got %v", expected, duration)
			}
		})

		t.Run("returns zero duration when StartedAt is nil", func(t *testing.T) {
			job := &Job{StartedAt: nil}
			currentTime := time.Date(2025, 11, 8, 10, 5, 30, 0, time.UTC)
			duration := job.GetExecutionDurationAt(currentTime)

			if duration != 0 {
				t.Errorf("expected duration 0, got %v", duration)
			}
		})

		t.Run("handles negative duration when currentTime is before StartedAt", func(t *testing.T) {
			startTime := time.Date(2025, 11, 8, 10, 5, 0, 0, time.UTC)
			job := &Job{StartedAt: &startTime}

			currentTime := time.Date(2025, 11, 8, 10, 0, 0, 0, time.UTC)
			duration := job.GetExecutionDurationAt(currentTime)

			expected := -5 * time.Minute
			if duration != expected {
				t.Errorf("expected duration %v, got %v", expected, duration)
			}
		})
	})
}
