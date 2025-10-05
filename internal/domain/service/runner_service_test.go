package service

import (
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

func TestUpdateRunnerStatus(t *testing.T) {
	service := NewRunnerService()

	t.Run("update runner status with active job", func(t *testing.T) {
		runners := []*entity.Runner{
			{ID: 1, Name: "runner-1", Status: entity.StatusIdle},
			{ID: 2, Name: "runner-2", Status: entity.StatusIdle},
		}

		runnerID := int64(1)
		jobs := []*entity.Job{
			{ID: 100, Name: "build", Status: "in_progress", RunnerID: &runnerID},
		}

		service.UpdateRunnerStatus(runners, jobs)

		if runners[0].Status != entity.StatusActive {
			t.Errorf("expected runner-1 to be Active, got %s", runners[0].Status)
		}

		if runners[1].Status != entity.StatusIdle {
			t.Errorf("expected runner-2 to be Idle, got %s", runners[1].Status)
		}
	})

	t.Run("no active jobs", func(t *testing.T) {
		runners := []*entity.Runner{
			{ID: 1, Name: "runner-1", Status: entity.StatusIdle},
		}

		jobs := []*entity.Job{}

		service.UpdateRunnerStatus(runners, jobs)

		if runners[0].Status != entity.StatusIdle {
			t.Errorf("expected runner-1 to remain Idle, got %s", runners[0].Status)
		}
	})

	t.Run("offline runner with job", func(t *testing.T) {
		runners := []*entity.Runner{
			{ID: 1, Name: "runner-1", Status: entity.StatusOffline},
		}

		runnerID := int64(1)
		jobs := []*entity.Job{
			{ID: 100, Name: "build", Status: "in_progress", RunnerID: &runnerID},
		}

		service.UpdateRunnerStatus(runners, jobs)

		// Offline runner should become active if it has a job
		if runners[0].Status != entity.StatusActive {
			t.Errorf("expected runner-1 to be Active, got %s", runners[0].Status)
		}
	})
}

func TestFindJobForRunner(t *testing.T) {
	service := NewRunnerService()
	startedAt := time.Now()

	t.Run("find job for runner", func(t *testing.T) {
		runnerID1 := int64(1)
		runnerID2 := int64(2)

		jobs := []*entity.Job{
			{ID: 100, Name: "build", Status: "in_progress", RunnerID: &runnerID1, StartedAt: &startedAt},
			{ID: 101, Name: "test", Status: "in_progress", RunnerID: &runnerID2, StartedAt: &startedAt},
		}

		job := service.FindJobForRunner(1, jobs)
		if job == nil {
			t.Fatal("expected to find job for runner 1")
		}

		if job.ID != 100 {
			t.Errorf("expected job ID 100, got %d", job.ID)
		}
	})

	t.Run("no job for runner", func(t *testing.T) {
		runnerID := int64(1)
		jobs := []*entity.Job{
			{ID: 100, Name: "build", Status: "in_progress", RunnerID: &runnerID, StartedAt: &startedAt},
		}

		job := service.FindJobForRunner(999, jobs)
		if job != nil {
			t.Errorf("expected no job for runner 999, got job %d", job.ID)
		}
	})

	t.Run("empty jobs list", func(t *testing.T) {
		jobs := []*entity.Job{}

		job := service.FindJobForRunner(1, jobs)
		if job != nil {
			t.Error("expected no job when jobs list is empty")
		}
	})
}
