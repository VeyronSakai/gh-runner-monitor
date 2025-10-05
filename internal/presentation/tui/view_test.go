package tui

import (
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   entity.RunnerStatus
		expected string
	}{
		{
			name:     "idle status",
			status:   entity.StatusIdle,
			expected: "🟢",
		},
		{
			name:     "active status",
			status:   entity.StatusActive,
			expected: "🟠",
		},
		{
			name:     "offline status",
			status:   entity.StatusOffline,
			expected: "⚫",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusIcon(tt.status)
			if result != tt.expected {
				t.Errorf("expected icon %s, got %s", tt.expected, result)
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
			name:     "less than one minute",
			duration: 45 * time.Second,
			expected: "00:45",
		},
		{
			name:     "one minute",
			duration: 1 * time.Minute,
			expected: "01:00",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "05:30",
		},
		{
			name:     "one hour",
			duration: 1 * time.Hour,
			expected: "01:00:00",
		},
		{
			name:     "hours, minutes, and seconds",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			expected: "02:15:45",
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
	t.Run("view with error", func(t *testing.T) {
		model := &Model{
			owner: "test-owner",
			repo:  "test-repo",
			err:   nil,
		}

		view := model.View()
		if view == "" {
			t.Error("expected non-empty view")
		}
	})

	t.Run("view when quitting", func(t *testing.T) {
		model := &Model{
			quitting: true,
		}

		view := model.View()
		if view != "" {
			t.Error("expected empty view when quitting")
		}
	})
}
