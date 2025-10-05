package entity

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

func TestRunnerMethods(t *testing.T) {
	tests := []struct {
		name      string
		status    RunnerStatus
		isOnline  bool
		isIdle    bool
		isActive  bool
	}{
		{
			name:     "idle runner",
			status:   StatusIdle,
			isOnline: true,
			isIdle:   true,
			isActive: false,
		},
		{
			name:     "active runner",
			status:   StatusActive,
			isOnline: true,
			isIdle:   false,
			isActive: true,
		},
		{
			name:     "offline runner",
			status:   StatusOffline,
			isOnline: false,
			isIdle:   false,
			isActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &Runner{Status: tt.status}

			if runner.IsOnline() != tt.isOnline {
				t.Errorf("IsOnline() = %v, want %v", runner.IsOnline(), tt.isOnline)
			}

			if runner.IsIdle() != tt.isIdle {
				t.Errorf("IsIdle() = %v, want %v", runner.IsIdle(), tt.isIdle)
			}

			if runner.IsActive() != tt.isActive {
				t.Errorf("IsActive() = %v, want %v", runner.IsActive(), tt.isActive)
			}
		})
	}
}
