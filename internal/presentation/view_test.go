package presentation

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

func TestUpdateTableRows_MatrixJobs(t *testing.T) {
	t.Run("matrix job with RunnerName but no RunnerID", func(t *testing.T) {
		// Setup test data simulating matrix job scenario
		runnerName := "test-runner"
		startedAt := time.Now().Add(-5 * time.Minute)

		runners := []*entity.Runner{
			{
				ID:     1,
				Name:   "test-runner",
				Status: entity.StatusActive,
			},
		}

		jobs := []*entity.Job{
			{
				ID:           101,
				RunID:        1001,
				Name:         "matrix-test (job-a)",
				Status:       "in_progress",
				RunnerID:     nil, // RunnerID is not set yet (matrix job scenario)
				RunnerName:   &runnerName,
				StartedAt:    &startedAt,
				WorkflowName: "Matrix CI",
			},
		}

		model := NewModel(nil, "owner", "repo", "", 5)
		model.runners = runners
		model.jobs = jobs
		model.currentTime = time.Now()

		// Call updateTableRows
		model.updateTableRows()

		// Verify that the job name is displayed (not "-")
		rows := model.table.Rows()
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}

		jobName := rows[0][2] // Job Name is the 3rd column
		if jobName == "-" {
			t.Error("expected job name to be displayed, but got '-'")
		}
		if jobName != "matrix-test (job-a) (Matrix CI)" {
			t.Errorf("expected job name 'matrix-test (job-a) (Matrix CI)', got '%s'", jobName)
		}

		execTime := rows[0][3] // Execution Time is the 4th column
		if execTime == "-" {
			t.Error("expected execution time to be displayed, but got '-'")
		}
	})

	t.Run("job with RunnerName matches by name", func(t *testing.T) {
		// Setup test data with both RunnerID and RunnerName set
		// Matching should be done by name only, not by ID
		runnerID := int64(999) // Different ID to ensure matching is by name only
		runnerName := "test-runner"
		startedAt := time.Now().Add(-3 * time.Minute)

		runners := []*entity.Runner{
			{
				ID:     1,
				Name:   "test-runner",
				Status: entity.StatusActive,
			},
		}

		jobs := []*entity.Job{
			{
				ID:           102,
				RunID:        1002,
				Name:         "normal-job",
				Status:       "in_progress",
				RunnerID:     &runnerID, // ID doesn't match, but name does
				RunnerName:   &runnerName,
				StartedAt:    &startedAt,
				WorkflowName: "CI",
			},
		}

		model := NewModel(nil, "owner", "repo", "", 5)
		model.runners = runners
		model.jobs = jobs
		model.currentTime = time.Now()

		// Call updateTableRows
		model.updateTableRows()

		// Verify that the job name is displayed (matched by name, not ID)
		rows := model.table.Rows()
		if len(rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rows))
		}

		jobName := rows[0][2]
		if jobName != "normal-job (CI)" {
			t.Errorf("expected job name 'normal-job (CI)', got '%s'", jobName)
		}
	})
}
