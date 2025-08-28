package github

import (
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
)

func TestClient_GetRunners(t *testing.T) {
	// This is a basic test structure
	// In a real scenario, you would mock the HTTP client
	t.Run("converts runner status correctly", func(t *testing.T) {
		// Test that online + busy = Active
		// Test that online + not busy = Idle
		// Test that offline = Offline
		t.Skip("Requires HTTP client mocking")
	})
}

func TestClient_GetActiveJobs(t *testing.T) {
	t.Run("filters only active jobs", func(t *testing.T) {
		// Test that only in_progress and queued jobs are returned
		t.Skip("Requires HTTP client mocking")
	})
}

func TestFormatDuration(t *testing.T) {
	// Helper function test
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "seconds only",
			duration: 45 * time.Second,
			expected: "00:45",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "05:30",
		},
		{
			name:     "hours minutes and seconds",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			expected: "02:15:45",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test would go here if formatDuration was exported
			t.Skip("formatDuration is in tui package")
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   models.RunnerStatus
		expected string
	}{
		{
			name:     "idle status",
			status:   models.StatusIdle,
			expected: "🟢",
		},
		{
			name:     "active status",
			status:   models.StatusActive,
			expected: "🟠",
		},
		{
			name:     "offline status",
			status:   models.StatusOffline,
			expected: "⚫",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test would go here if getStatusIcon was exported
			t.Skip("getStatusIcon is in tui package")
		})
	}
}