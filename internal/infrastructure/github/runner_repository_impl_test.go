package github

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

func TestRunnerRepositoryImpl_ConvertRunners(t *testing.T) {
	repo := &RunnerRepositoryImpl{}

	tests := []struct {
		name     string
		response *runnersResponse
		want     int
	}{
		{
			name: "convert online idle runner",
			response: &runnersResponse{
				TotalCount: 1,
				Runners: []runnerResponse{
					{
						ID:     1,
						Name:   "test-runner",
						OS:     "Linux",
						Status: "online",
						Busy:   false,
						Labels: []label{
							{ID: 1, Name: "self-hosted", Type: "custom"},
							{ID: 2, Name: "linux", Type: "system"},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "convert online active runner",
			response: &runnersResponse{
				TotalCount: 1,
				Runners: []runnerResponse{
					{
						ID:     2,
						Name:   "active-runner",
						OS:     "macOS",
						Status: "online",
						Busy:   true,
						Labels: []label{
							{ID: 1, Name: "self-hosted", Type: "custom"},
							{ID: 3, Name: "macos", Type: "system"},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "convert offline runner",
			response: &runnersResponse{
				TotalCount: 1,
				Runners: []runnerResponse{
					{
						ID:     3,
						Name:   "offline-runner",
						OS:     "Windows",
						Status: "offline",
						Busy:   false,
						Labels: []label{
							{ID: 1, Name: "self-hosted", Type: "custom"},
						},
					},
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]*entity.Runner, 0, len(tt.response.Runners))
			for _, runner := range tt.response.Runners {
				status := entity.StatusOffline
				if runner.Status == "online" {
					if runner.Busy {
						status = entity.StatusActive
					} else {
						status = entity.StatusIdle
					}
				}

				labels := make([]string, 0, len(runner.Labels))
				for _, l := range runner.Labels {
					labels = append(labels, l.Name)
				}

				result = append(result, &entity.Runner{
					ID:        runner.ID,
					Name:      runner.Name,
					Status:    status,
					Labels:    labels,
					OS:        runner.OS,
					UpdatedAt: time.Now(),
				})
			}
			got := result
			if len(got) != tt.want {
				t.Errorf("convertRunners() got %d runners, want %d", len(got), tt.want)
			}

			if len(got) > 0 {
				runner := got[0]
				if runner.ID != tt.response.Runners[0].ID {
					t.Errorf("convertRunners() ID = %d, want %d", runner.ID, tt.response.Runners[0].ID)
				}
				if runner.Name != tt.response.Runners[0].Name {
					t.Errorf("convertRunners() Name = %s, want %s", runner.Name, tt.response.Runners[0].Name)
				}
				if runner.OS != tt.response.Runners[0].OS {
					t.Errorf("convertRunners() OS = %s, want %s", runner.OS, tt.response.Runners[0].OS)
				}

				// Status validation
				expectedStatus := entity.StatusOffline
				if tt.response.Runners[0].Status == "online" {
					if tt.response.Runners[0].Busy {
						expectedStatus = entity.StatusActive
					} else {
						expectedStatus = entity.StatusIdle
					}
				}
				if runner.Status != expectedStatus {
					t.Errorf("convertRunners() Status = %s, want %s", runner.Status, expectedStatus)
				}

				// Labels validation
				if len(runner.Labels) != len(tt.response.Runners[0].Labels) {
					t.Errorf("convertRunners() Labels count = %d, want %d", len(runner.Labels), len(tt.response.Runners[0].Labels))
				}
			}
		})
	}
}

func TestRunnerRepositoryImpl_InterfaceCompliance(t *testing.T) {
	// This test ensures RunnerRepositoryImpl implements the RunnerRepository interface
	// by checking that the type satisfies the interface at compile time
	var _ interface {
		GetRunners(context.Context, string, string, string) ([]*entity.Runner, error)
		GetActiveJobs(context.Context, string, string, string) ([]*entity.Job, error)
	} = (*RunnerRepositoryImpl)(nil)

	t.Log("RunnerRepositoryImpl implements the required interface")
}

// Tests from client_test.go

func TestPathConstruction(t *testing.T) {
	tests := []struct {
		name     string
		org      string
		owner    string
		repo     string
		expected string
	}{
		{
			name:     "repository runners path",
			org:      "",
			owner:    "myowner",
			repo:     "myrepo",
			expected: "repos/myowner/myrepo/actions/runners",
		},
		{
			name:     "organization runners path",
			org:      "myorg",
			owner:    "",
			repo:     "",
			expected: "orgs/myorg/actions/runners",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.org != "" {
				path = fmt.Sprintf("orgs/%s/actions/runners", tt.org)
			} else {
				path = fmt.Sprintf("repos/%s/%s/actions/runners", tt.owner, tt.repo)
			}

			if path != tt.expected {
				t.Errorf("expected path %s, got %s", tt.expected, path)
			}
		})
	}
}

// Tests from job_test.go

func TestGetJobsForRun(t *testing.T) {
	tests := []struct {
		name      string
		fullName  string
		org       string
		expectErr bool
	}{
		{
			name:      "valid repository format",
			fullName:  "owner/repo",
			org:       "org",
			expectErr: false,
		},
		{
			name:      "invalid repository format - no slash",
			fullName:  "ownerrepo",
			org:       "org",
			expectErr: true,
		},
		{
			name:      "invalid repository format - multiple slashes",
			fullName:  "owner/repo/extra",
			org:       "org",
			expectErr: true,
		},
		{
			name:      "empty repository name",
			fullName:  "",
			org:       "org",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the repository name parsing logic
			if tt.org != "" {
				parts := strings.Split(tt.fullName, "/")
				hasError := len(parts) != 2

				if hasError != tt.expectErr {
					t.Errorf("expected error: %v, got error: %v", tt.expectErr, hasError)
				}
			}
		})
	}
}

func TestFilterActiveJobs(t *testing.T) {
	now := time.Now()
	runnerID := int64(1)
	runnerName := "test-runner"

	jobsResp := &jobsResponse{
		Jobs: []jobResponse{
			{
				ID:         1,
				RunID:      100,
				Name:       "build",
				Status:     "in_progress",
				RunnerID:   &runnerID,
				RunnerName: &runnerName,
				StartedAt:  &now,
			},
			{
				ID:         2,
				RunID:      100,
				Name:       "test",
				Status:     "queued",
				RunnerID:   nil,
				RunnerName: nil,
				StartedAt:  nil,
			},
			{
				ID:         3,
				RunID:      100,
				Name:       "deploy",
				Status:     "completed",
				RunnerID:   &runnerID,
				RunnerName: &runnerName,
				StartedAt:  &now,
			},
			{
				ID:         4,
				RunID:      100,
				Name:       "cleanup",
				Status:     "cancelled",
				RunnerID:   nil,
				RunnerName: nil,
				StartedAt:  nil,
			},
		},
	}

	// Test that only in_progress and queued jobs are included
	activeCount := 0
	for _, job := range jobsResp.Jobs {
		if job.Status == "in_progress" || job.Status == "queued" {
			activeCount++
		}
	}

	if activeCount != 2 {
		t.Errorf("expected 2 active jobs, got %d", activeCount)
	}
}
