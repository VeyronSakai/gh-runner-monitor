package tui

import (
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
)

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   models.RunnerStatus
		expected string
	}{
		{
			name:     "idle status shows green circle",
			status:   models.StatusIdle,
			expected: "🟢",
		},
		{
			name:     "active status shows orange circle",
			status:   models.StatusActive,
			expected: "🟠",
		},
		{
			name:     "offline status shows black circle",
			status:   models.StatusOffline,
			expected: "⚫",
		},
		{
			name:     "unknown status shows question mark",
			status:   models.RunnerStatus("unknown"),
			expected: "❓",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusIcon(tt.status)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
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
			name:     "exactly one minute",
			duration: 60 * time.Second,
			expected: "01:00",
		},
		{
			name:     "hours minutes and seconds",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			expected: "02:15:45",
		},
		{
			name:     "exactly one hour",
			duration: 3600 * time.Second,
			expected: "01:00:00",
		},
		{
			name:     "over 24 hours",
			duration: 25*time.Hour + 30*time.Minute,
			expected: "25:30:00",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestView(t *testing.T) {
	// Note: This test would require mocking the github.Client
	// For now, we'll skip the actual view testing
	t.Run("view structure test", func(t *testing.T) {
		t.Skip("Requires github.Client mocking")
	})
}