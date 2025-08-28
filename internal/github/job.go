package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
)

// GetActiveJobs fetches active jobs for a repository or organization
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
		jobs, err := c.getJobsForRun(run, org, owner, repo)
		if err != nil {
			continue // Skip this run if we can't get jobs
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
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

func (c *Client) getJobsForRun(run workflowRun, org, owner, repo string) ([]*models.Job, error) {
	// For org-level, we need to extract owner and repo from the full name
	var runOwner, runRepo string
	if org != "" {
		// Parse repository full name (format: "owner/repo")
		parts := strings.Split(run.Repository.FullName, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid repository full name format: %s", run.Repository.FullName)
		}
		runOwner = parts[0]
		runRepo = parts[1]
	} else {
		runOwner = owner
		runRepo = repo
	}

	path := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", runOwner, runRepo, run.ID)
	jobs, err := c.fetchJobs(path)
	if err != nil {
		return nil, err
	}

	var result []*models.Job
	for _, job := range jobs.Jobs {
		if job.Status == "in_progress" || job.Status == "queued" {
			result = append(result, &models.Job{
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

	return result, nil
}
