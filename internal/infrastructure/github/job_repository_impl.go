package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	domainrepo "github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/cli/go-gh/v2/pkg/api"
)

// JobRepositoryImpl implements the JobRepository interface using GitHub API
type JobRepositoryImpl struct {
	restClient *api.RESTClient
}

// NewJobRepository creates a new instance of JobRepositoryImpl
func NewJobRepository() (domainrepo.JobRepository, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w\nPlease run 'gh auth login' to authenticate with GitHub", err)
	}

	return &JobRepositoryImpl{
		restClient: restClient,
	}, nil
}

// FetchActiveJobs retrieves all active jobs for a repository or organization
func (j *JobRepositoryImpl) FetchActiveJobs(ctx context.Context, owner, repo, org string) ([]*entity.Job, error) {
	var allJobs []*entity.Job

	// Fetch in_progress workflow runs
	inProgressPath := j.getWorkflowRunsPath(owner, repo, org, "in_progress")
	inProgressRuns, err := j.fetchWorkflowRuns(inProgressPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch in_progress runs: %w", err)
	}

	for _, run := range inProgressRuns.WorkflowRuns {
		jobs, err := j.getJobsForRun(run, org, owner, repo)
		if err != nil {
			continue // Skip this run if we can't get jobs
		}
		allJobs = append(allJobs, jobs...)
	}

	// Fetch queued workflow runs
	queuedPath := j.getWorkflowRunsPath(owner, repo, org, "queued")
	queuedRuns, err := j.fetchWorkflowRuns(queuedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch queued runs: %w", err)
	}

	for _, run := range queuedRuns.WorkflowRuns {
		jobs, err := j.getJobsForRun(run, org, owner, repo)
		if err != nil {
			continue // Skip this run if we can't get jobs
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

// getWorkflowRunsPath constructs the API path for fetching workflow runs with a specific status
func (j *JobRepositoryImpl) getWorkflowRunsPath(owner, repo, org, status string) string {
	if org != "" {
		return fmt.Sprintf("orgs/%s/actions/runs?status=%s", org, status)
	}
	return fmt.Sprintf("repos/%s/%s/actions/runs?status=%s", owner, repo, status)
}

// fetchWorkflowRuns fetches workflow runs from GitHub API with pagination support
func (j *JobRepositoryImpl) fetchWorkflowRuns(path string) (*workflowRunsResponse, error) {
	allRuns := &workflowRunsResponse{
		WorkflowRuns: []workflowRun{},
	}

	// Determine the separator for query parameters
	// Use "?" if the path doesn't have any query parameters yet
	// Use "&" if the path already has query parameters (e.g., "?status=in_progress")
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}

	page := 1
	perPage := 100

	for {
		currentPath := fmt.Sprintf("%s%sper_page=%d&page=%d", path, separator, perPage, page)
		response, err := j.restClient.Request(http.MethodGet, currentPath, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to request workflow runs: %w", err)
		}

		var runs workflowRunsResponse
		if err := json.NewDecoder(response.Body).Decode(&runs); err != nil {
			_ = response.Body.Close()
			return nil, fmt.Errorf("failed to decode workflow runs response: %w", err)
		}
		_ = response.Body.Close()

		allRuns.TotalCount = runs.TotalCount
		allRuns.WorkflowRuns = append(allRuns.WorkflowRuns, runs.WorkflowRuns...)

		// If we got fewer items than per_page, we've reached the last page
		if len(runs.WorkflowRuns) < perPage {
			break
		}

		page++
	}

	return allRuns, nil
}

// requestGetJobs fetches jobs from GitHub API
func (j *JobRepositoryImpl) requestGetJobs(path string) (*jobsResponse, error) {
	response, err := j.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request jobs: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var jobs jobsResponse
	if err := json.NewDecoder(response.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs response: %w", err)
	}

	return &jobs, nil
}

// getJobsForRun fetches and converts jobs for a specific workflow run
func (j *JobRepositoryImpl) getJobsForRun(run workflowRun, org, owner, repo string) ([]*entity.Job, error) {
	runOwner, runRepo, err := j.extractOwnerAndRepo(org, owner, repo, run.Repository.FullName)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", runOwner, runRepo, run.ID)
	jobs, err := j.requestGetJobs(path)
	if err != nil {
		return nil, err
	}

	var result []*entity.Job
	for _, job := range jobs.Jobs {
		if j.isActiveJob(job.Status) {
			result = append(result, &entity.Job{
				ID:           job.ID,
				RunID:        job.RunID,
				Name:         job.Name,
				Status:       job.Status,
				RunnerID:     job.RunnerID,
				RunnerName:   job.RunnerName,
				StartedAt:    job.StartedAt,
				WorkflowName: run.Name,
				Repository:   run.Repository.FullName,
				HtmlUrl:      job.HtmlUrl,
			})
		}
	}

	return result, nil
}

// extractOwnerAndRepo extracts owner and repo from either org context or direct parameters
func (j *JobRepositoryImpl) extractOwnerAndRepo(org, owner, repo, fullName string) (string, string, error) {
	if org != "" {
		// Parse repository full name (format: "owner/repo")
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid repository full name format: %s", fullName)
		}
		return parts[0], parts[1], nil
	}
	return owner, repo, nil
}

// isActiveJob checks if a job status is considered active (in_progress or queued)
func (j *JobRepositoryImpl) isActiveJob(status string) bool {
	return status == "in_progress" || status == "queued"
}
