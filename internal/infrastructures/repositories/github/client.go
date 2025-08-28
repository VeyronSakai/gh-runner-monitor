package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh/v2/pkg/api"
	"net/http"
	"time"
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

type RunnerResponse struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	OS        string   `json:"os"`
	Status    string   `json:"status"`
	Busy      bool     `json:"busy"`
	Labels    []Label  `json:"labels"`
}

type Label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RunnersResponse struct {
	TotalCount int              `json:"total_count"`
	Runners    []RunnerResponse `json:"runners"`
}

type WorkflowRunsResponse struct {
	TotalCount   int            `json:"total_count"`
	WorkflowRuns []WorkflowRun  `json:"workflow_runs"`
}

type WorkflowRun struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Conclusion *string    `json:"conclusion"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type JobsResponse struct {
	TotalCount int           `json:"total_count"`
	Jobs       []JobResponse `json:"jobs"`
}

type JobResponse struct {
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

func (c *Client) GetRepoRunners(ctx context.Context, owner, repo string) (*RunnersResponse, error) {
	path := fmt.Sprintf("repos/%s/%s/actions/runners", owner, repo)
	
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request runners: %w", err)
	}
	defer response.Body.Close()
	
	var runners RunnersResponse
	if err := json.NewDecoder(response.Body).Decode(&runners); err != nil {
		return nil, fmt.Errorf("failed to decode runners response: %w", err)
	}
	
	return &runners, nil
}

func (c *Client) GetOrgRunners(ctx context.Context, org string) (*RunnersResponse, error) {
	path := fmt.Sprintf("orgs/%s/actions/runners", org)
	
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request org runners: %w", err)
	}
	defer response.Body.Close()
	
	var runners RunnersResponse
	if err := json.NewDecoder(response.Body).Decode(&runners); err != nil {
		return nil, fmt.Errorf("failed to decode org runners response: %w", err)
	}
	
	return &runners, nil
}

func (c *Client) GetRepoActiveWorkflowRuns(ctx context.Context, owner, repo string) (*WorkflowRunsResponse, error) {
	path := fmt.Sprintf("repos/%s/%s/actions/runs?status=in_progress", owner, repo)
	
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request workflow runs: %w", err)
	}
	defer response.Body.Close()
	
	var runs WorkflowRunsResponse
	if err := json.NewDecoder(response.Body).Decode(&runs); err != nil {
		return nil, fmt.Errorf("failed to decode workflow runs response: %w", err)
	}
	
	return &runs, nil
}

func (c *Client) GetOrgActiveWorkflowRuns(ctx context.Context, org string) (*WorkflowRunsResponse, error) {
	path := fmt.Sprintf("orgs/%s/actions/runs?status=in_progress", org)
	
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request org workflow runs: %w", err)
	}
	defer response.Body.Close()
	
	var runs WorkflowRunsResponse
	if err := json.NewDecoder(response.Body).Decode(&runs); err != nil {
		return nil, fmt.Errorf("failed to decode org workflow runs response: %w", err)
	}
	
	return &runs, nil
}

func (c *Client) GetWorkflowRunJobs(ctx context.Context, owner, repo string, runID int64) (*JobsResponse, error) {
	path := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID)
	
	response, err := c.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request jobs: %w", err)
	}
	defer response.Body.Close()
	
	var jobs JobsResponse
	if err := json.NewDecoder(response.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs response: %w", err)
	}
	
	return &jobs, nil
}