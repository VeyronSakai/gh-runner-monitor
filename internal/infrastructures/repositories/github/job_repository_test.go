package github

import (
	"testing"
	"time"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

func TestJobRepository_convertToEntity(t *testing.T) {
	repo := &JobRepository{}
	
	t.Run("convert in_progress job", func(t *testing.T) {
		runnerID := int64(1)
		runnerName := "test-runner"
		startedAt := time.Now().Add(-10 * time.Minute)
		
		resp := JobResponse{
			ID:         100,
			RunID:      200,
			Name:       "Test Job",
			Status:     "in_progress",
			RunnerID:   &runnerID,
			RunnerName: &runnerName,
			StartedAt:  &startedAt,
		}
		
		job := repo.convertToEntity(resp, "Test Workflow", "owner/repo")
		
		if job.ID != 100 {
			t.Errorf("expected ID 100, got %d", job.ID)
		}
		if job.RunID != 200 {
			t.Errorf("expected RunID 200, got %d", job.RunID)
		}
		if job.Name != "Test Job" {
			t.Errorf("expected name Test Job, got %s", job.Name)
		}
		if job.Status != entities.JobStatusInProgress {
			t.Errorf("expected status InProgress, got %s", job.Status)
		}
		if !job.IsRunning() {
			t.Error("expected job to be running")
		}
		if job.GetRunnerID() != 1 {
			t.Errorf("expected runner ID 1, got %d", job.GetRunnerID())
		}
		if job.WorkflowName != "Test Workflow" {
			t.Errorf("expected workflow name Test Workflow, got %s", job.WorkflowName)
		}
		if job.Repository != "owner/repo" {
			t.Errorf("expected repository owner/repo, got %s", job.Repository)
		}
		if job.GetExecutionTime() == 0 {
			t.Error("expected non-zero execution time")
		}
	})
	
	t.Run("convert queued job", func(t *testing.T) {
		resp := JobResponse{
			ID:     101,
			RunID:  201,
			Name:   "Queued Job",
			Status: "queued",
		}
		
		job := repo.convertToEntity(resp, "Workflow", "owner/repo")
		
		if job.Status != entities.JobStatusQueued {
			t.Errorf("expected status Queued, got %s", job.Status)
		}
		if job.IsRunning() {
			t.Error("expected job not to be running")
		}
		if job.GetRunnerID() != 0 {
			t.Errorf("expected runner ID 0, got %d", job.GetRunnerID())
		}
		if job.GetExecutionTime() != 0 {
			t.Error("expected zero execution time for queued job")
		}
	})
	
	t.Run("convert completed job", func(t *testing.T) {
		conclusion := "success"
		startedAt := time.Now().Add(-20 * time.Minute)
		completedAt := time.Now().Add(-5 * time.Minute)
		
		resp := JobResponse{
			ID:          102,
			RunID:       202,
			Name:        "Completed Job",
			Status:      "completed",
			Conclusion:  &conclusion,
			StartedAt:   &startedAt,
			CompletedAt: &completedAt,
		}
		
		job := repo.convertToEntity(resp, "Workflow", "owner/repo")
		
		if job.Status != entities.JobStatusCompleted {
			t.Errorf("expected status Completed, got %s", job.Status)
		}
		if job.IsRunning() {
			t.Error("expected job not to be running")
		}
		if job.Conclusion != "success" {
			t.Errorf("expected conclusion success, got %s", job.Conclusion)
		}
		
		expectedDuration := completedAt.Sub(startedAt)
		actualDuration := job.GetExecutionTime()
		if actualDuration != expectedDuration {
			t.Errorf("expected execution time %v, got %v", expectedDuration, actualDuration)
		}
	})
}