package github

import (
	"testing"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

func TestRunnerRepository_convertToEntity(t *testing.T) {
	repo := &RunnerRepository{}
	
	t.Run("convert offline runner", func(t *testing.T) {
		resp := RunnerResponse{
			ID:     1,
			Name:   "test-runner",
			OS:     "linux",
			Status: "offline",
			Busy:   false,
			Labels: []Label{
				{ID: 1, Name: "self-hosted", Type: "read-only"},
				{ID: 2, Name: "linux", Type: "read-only"},
			},
		}
		
		runner := repo.convertToEntity(resp)
		
		if runner.ID != 1 {
			t.Errorf("expected ID 1, got %d", runner.ID)
		}
		if runner.Name != "test-runner" {
			t.Errorf("expected name test-runner, got %s", runner.Name)
		}
		if runner.Status != entities.StatusOffline {
			t.Errorf("expected status Offline, got %s", runner.Status)
		}
		if len(runner.Labels) != 2 {
			t.Errorf("expected 2 labels, got %d", len(runner.Labels))
		}
		if runner.OS != "linux" {
			t.Errorf("expected OS linux, got %s", runner.OS)
		}
	})
	
	t.Run("convert busy runner", func(t *testing.T) {
		resp := RunnerResponse{
			ID:     2,
			Name:   "busy-runner",
			OS:     "windows",
			Status: "online",
			Busy:   true,
			Labels: []Label{
				{ID: 1, Name: "self-hosted", Type: "read-only"},
			},
		}
		
		runner := repo.convertToEntity(resp)
		
		if runner.Status != entities.StatusActive {
			t.Errorf("expected status Active, got %s", runner.Status)
		}
	})
	
	t.Run("convert idle runner", func(t *testing.T) {
		resp := RunnerResponse{
			ID:     3,
			Name:   "idle-runner",
			OS:     "macos",
			Status: "online",
			Busy:   false,
			Labels: []Label{},
		}
		
		runner := repo.convertToEntity(resp)
		
		if runner.Status != entities.StatusIdle {
			t.Errorf("expected status Idle, got %s", runner.Status)
		}
		if len(runner.Labels) != 0 {
			t.Errorf("expected 0 labels, got %d", len(runner.Labels))
		}
	})
}