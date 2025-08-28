package github

import (
	"testing"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
)

func TestConvertRunners(t *testing.T) {
	client := &Client{}
	
	tests := []struct {
		name     string
		input    *runnersResponse
		expected int
	}{
		{
			name: "convert multiple runners",
			input: &runnersResponse{
				TotalCount: 3,
				Runners: []runnerResponse{
					{
						ID:     1,
						Name:   "runner-1",
						Status: "online",
						Busy:   false,
						OS:     "linux",
						Labels: []label{{Name: "self-hosted"}},
					},
					{
						ID:     2,
						Name:   "runner-2",
						Status: "online",
						Busy:   true,
						OS:     "windows",
						Labels: []label{{Name: "self-hosted"}, {Name: "windows"}},
					},
					{
						ID:     3,
						Name:   "runner-3",
						Status: "offline",
						Busy:   false,
						OS:     "macos",
						Labels: []label{{Name: "self-hosted"}, {Name: "macos"}},
					},
				},
			},
			expected: 3,
		},
		{
			name: "handle empty response",
			input: &runnersResponse{
				TotalCount: 0,
				Runners:    []runnerResponse{},
			},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.convertRunners(tt.input)
			
			if len(result) != tt.expected {
				t.Errorf("expected %d runners, got %d", tt.expected, len(result))
			}
			
			if tt.expected > 0 {
				// Test status conversion
				if result[0].Status != models.StatusIdle {
					t.Errorf("expected first runner to be Idle, got %s", result[0].Status)
				}
				
				if tt.expected > 1 && result[1].Status != models.StatusActive {
					t.Errorf("expected second runner to be Active, got %s", result[1].Status)
				}
				
				if tt.expected > 2 && result[2].Status != models.StatusOffline {
					t.Errorf("expected third runner to be Offline, got %s", result[2].Status)
				}
			}
		})
	}
}

func TestRunnerStatusConversion(t *testing.T) {
	client := &Client{}
	
	tests := []struct {
		name         string
		status       string
		busy         bool
		expectedStatus models.RunnerStatus
	}{
		{
			name:         "online and not busy should be idle",
			status:       "online",
			busy:         false,
			expectedStatus: models.StatusIdle,
		},
		{
			name:         "online and busy should be active",
			status:       "online",
			busy:         true,
			expectedStatus: models.StatusActive,
		},
		{
			name:         "offline should be offline",
			status:       "offline",
			busy:         false,
			expectedStatus: models.StatusOffline,
		},
		{
			name:         "offline and busy should still be offline",
			status:       "offline",
			busy:         true,
			expectedStatus: models.StatusOffline,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := &runnersResponse{
				Runners: []runnerResponse{
					{
						ID:     1,
						Name:   "test-runner",
						Status: tt.status,
						Busy:   tt.busy,
						OS:     "linux",
					},
				},
			}
			
			result := client.convertRunners(response)
			
			if len(result) != 1 {
				t.Fatalf("expected 1 runner, got %d", len(result))
			}
			
			if result[0].Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, result[0].Status)
			}
		})
	}
}