package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
	"github.com/cli/go-gh/v2/pkg/api"
)

type Client struct {
	restClient *api.RESTClient
}

func NewClient() (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}
	
	return &Client{
		restClient: restClient,
	}, nil
}

type runnersResponse struct {
	TotalCount int              `json:"total_count"`
	Runners    []runnerResponse `json:"runners"`
}

type runnerResponse struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	OS        string   `json:"os"`
	Status    string   `json:"status"`
	Busy      bool     `json:"busy"`
	Labels    []label  `json:"labels"`
}

type label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type workflowRunsResponse struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []workflowRun `json:"workflow_runs"`
}

type workflowRun struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Conclusion *string    `json:"conclusion"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Repository repository `json:"repository"`
}

type repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type jobsResponse struct {
	TotalCount int           `json:"total_count"`
	Jobs       []jobResponse `json:"jobs"`
}

type jobResponse struct {
	ID          int64      `json:"id"`
	RunID       int64      `json:"run_id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  *string    `json:"conclusion"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	RunnerID    *int64     `json:"runner_id"`
	RunnerName  *string    `json:"runner_name"`
}

func (c *Client) GetRunners(ctx context.Context, owner, repo, org string) ([]*models.Runner, error) {
	var runners *runnersResponse
	var err error
	
	if org != "" {
		path := fmt.Sprintf("orgs/%s/actions/runners", org)
		runners, err = c.fetchRunners(path)
	} else {
		path := fmt.Sprintf("repos/%s/%s/actions/runners", owner, repo)
		runners, err = c.fetchRunners(path)
	}
	
	if err != nil {
		return nil, err
	}
	
	result := make([]*models.Runner, 0, len(runners.Runners))
	for _, r := range runners.Runners {
		status := models.StatusOffline
		if r.Status == "online" {
			if r.Busy {
				status = models.StatusActive
			} else {
				status = models.StatusIdle
			}
		}
		
		labels := make([]string, 0, len(r.Labels))
		for _, l := range r.Labels {
			labels = append(labels, l.Name)
		}
		
		result = append(result, &models.Runner{
			ID:        r.ID,
			Name:      r.Name,
			Status:    status,
			Labels:    labels,
			OS:        r.OS,
			UpdatedAt: time.Now(),
		})
	}
	
	return result, nil
}

func (c *Client) GetActiveJobs(ctx context.Context, owner, repo, org string) ([]*models.Job, error) {
	var runs *workflowRunsResponse
	var err error
	
	if org != "" {
		path := fmt.Sprintf("orgs/%s/actions/runs?status=in_progress", org)
		runs, err = c.fetchWorkflowRuns(path)
	} else {
		path := fmt.Sprintf("repos/%s/%s/actions/runs?status=in_progress", owner, repo)
		runs, err = c.fetchWorkflowRuns(path)
	}
	
	if err != nil {
		return nil, err
	}
	
	var allJobs []*models.Job
	for _, run := range runs.WorkflowRuns {
		// For org-level, we need to extract owner and repo from the full name
		var runOwner, runRepo string
		if org != "" {
			// Parse repository full name (format: "owner/repo")
			repoFullName := run.Repository.FullName
			var repoOwner, repoName string
			fmt.Sscanf(repoFullName, "%[^/]/%s", &repoOwner, &repoName)
			runOwner = repoOwner
			runRepo = repoName
		} else {
			runOwner = owner
			runRepo = repo
		}
		
		path := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", runOwner, runRepo, run.ID)
		jobs, err := c.fetchJobs(path)
		if err != nil {
			continue // Skip this run if we can't get jobs
		}
		
		for _, job := range jobs.Jobs {
			if job.Status == "in_progress" || job.Status == "queued" {
				allJobs = append(allJobs, &models.Job{
					ID:           job.ID,
					RunID:        job.RunID,
					Name:         job.Name,
					Status:       job.Status,
					RunnerID:     job.RunnerID,
					RunnerName:   job.RunnerName,
					StartedAt:    job.StartedAt,
					WorkflowName: run.Name,
					Repository:   run.Repository.FullName,
				})
			}
		}
	}
	
	return allJobs, nil
}

func (c *Client) fetchRunners(path string) (*runnersResponse, error) {
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request runners: %w", err)
	}
	defer response.Body.Close()
	
	var runners runnersResponse
	if err := json.NewDecoder(response.Body).Decode(&runners); err != nil {
		return nil, fmt.Errorf("failed to decode runners response: %w", err)
	}
	
	return &runners, nil
}

func (c *Client) fetchWorkflowRuns(path string) (*workflowRunsResponse, error) {
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request workflow runs: %w", err)
	}
	defer response.Body.Close()
	
	var runs workflowRunsResponse
	if err := json.NewDecoder(response.Body).Decode(&runs); err != nil {
		return nil, fmt.Errorf("failed to decode workflow runs response: %w", err)
	}
	
	return &runs, nil
}

func (c *Client) fetchJobs(path string) (*jobsResponse, error) {
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request jobs: %w", err)
	}
	defer response.Body.Close()
	
	var jobs jobsResponse
	if err := json.NewDecoder(response.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs response: %w", err)
	}
	
	return &jobs, nil
}