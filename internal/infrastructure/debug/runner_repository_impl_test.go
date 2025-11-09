package debug

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

func TestDebugRunnerRepository(t *testing.T) {
	// Create a temporary JSON file for testing
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "test_data.json")

	jsonContent := `{
  "CurrentTime": "2025-11-03T10:05:00Z",
  "runners": [
    {
      "ID": 1,
      "Name": "test-runner",
      "Status": "Active",
      "Labels": ["self-hosted", "linux"],
      "OS": "linux",
      "UpdatedAt": "2025-11-03T10:00:00Z"
    }
  ],
  "jobs": [
    {
      "ID": 1,
      "RunID": 100,
      "Name": "Test Job",
      "Status": "in_progress",
      "RunnerID": 1,
      "RunnerName": "test-runner",
      "StartedAt": "2025-11-03T10:00:30Z",
      "WorkflowName": "Test Workflow",
      "Repository": "test/repo"
    }
  ]
}`

	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	// Test creating the repository
	repo, err := NewRunnerRepository(jsonPath)
	if err != nil {
		t.Fatalf("Failed to create debug repository: %v", err)
	}

	ctx := context.Background()

	// Test FetchRunners
	runners, err := repo.FetchRunners(ctx, "", "", "")
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	if len(runners) != 1 {
		t.Errorf("Expected 1 runner, got %d", len(runners))
	}

	if runners[0].Name != "test-runner" {
		t.Errorf("Expected runner name 'test-runner', got '%s'", runners[0].Name)
	}

	if runners[0].Status != entity.StatusActive {
		t.Errorf("Expected runner status Active, got %s", runners[0].Status)
	}

	// Test FetchActiveJobs
	jobs, err := repo.FetchActiveJobs(ctx, "", "", "")
	if err != nil {
		t.Fatalf("FetchActiveJobs failed: %v", err)
	}

	if len(jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(jobs))
	}

	if jobs[0].Name != "Test Job" {
		t.Errorf("Expected job name 'Test Job', got '%s'", jobs[0].Name)
	}

	if jobs[0].Status != "in_progress" {
		t.Errorf("Expected job status 'in_progress', got '%s'", jobs[0].Status)
	}

	// Test GetCurrentTime
	currentTime := repo.(*RunnerRepositoryImpl).GetCurrentTime()
	expectedTime, _ := time.Parse(time.RFC3339, "2025-11-03T10:05:00Z")
	if !currentTime.Equal(expectedTime) {
		t.Errorf("Expected current time %v, got %v", expectedTime, currentTime)
	}
}

func TestDebugRunnerRepository_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "invalid.json")

	invalidJSON := `{ "invalid": json }`

	if err := os.WriteFile(jsonPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}

	_, err := NewRunnerRepository(jsonPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestDebugRunnerRepository_FileNotFound(t *testing.T) {
	_, err := NewRunnerRepository("/nonexistent/path/data.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
