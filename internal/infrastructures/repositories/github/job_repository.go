package github

import (
	"context"
	"fmt"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
	"strings"
)

type JobRepository struct {
	client *Client
}

func NewJobRepository(client *Client) *JobRepository {
	return &JobRepository{
		client: client,
	}
}

func (r *JobRepository) ListActiveJobs(ctx context.Context, owner, repo string) ([]*entities.Job, error) {
	workflowRuns, err := r.client.GetRepoActiveWorkflowRuns(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get active workflow runs: %w", err)
	}
	
	allJobs := make([]*entities.Job, 0)
	
	for _, run := range workflowRuns.WorkflowRuns {
		jobsResp, err := r.client.GetWorkflowRunJobs(ctx, owner, repo, run.ID)
		if err != nil {
			continue
		}
		
		for _, jobResp := range jobsResp.Jobs {
			job := r.convertToEntity(jobResp, run.Name, fmt.Sprintf("%s/%s", owner, repo))
			if job.IsRunning() {
				allJobs = append(allJobs, job)
			}
		}
	}
	
	return allJobs, nil
}

func (r *JobRepository) ListOrgActiveJobs(ctx context.Context, org string) ([]*entities.Job, error) {
	workflowRuns, err := r.client.GetOrgActiveWorkflowRuns(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to get org active workflow runs: %w", err)
	}
	
	allJobs := make([]*entities.Job, 0)
	
	for _, run := range workflowRuns.WorkflowRuns {
		parts := strings.Split(run.Repository.FullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repo := parts[0], parts[1]
		
		jobsResp, err := r.client.GetWorkflowRunJobs(ctx, owner, repo, run.ID)
		if err != nil {
			continue
		}
		
		for _, jobResp := range jobsResp.Jobs {
			job := r.convertToEntity(jobResp, run.Name, run.Repository.FullName)
			if job.IsRunning() {
				allJobs = append(allJobs, job)
			}
		}
	}
	
	return allJobs, nil
}

func (r *JobRepository) GetJob(ctx context.Context, owner, repo string, jobID int64) (*entities.Job, error) {
	jobs, err := r.ListActiveJobs(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	
	for _, job := range jobs {
		if job.ID == jobID {
			return job, nil
		}
	}
	
	return nil, fmt.Errorf("job with ID %d not found", jobID)
}

func (r *JobRepository) convertToEntity(resp JobResponse, workflowName, repository string) *entities.Job {
	var status entities.JobStatus
	
	switch resp.Status {
	case "queued":
		status = entities.JobStatusQueued
	case "in_progress":
		status = entities.JobStatusInProgress
	case "completed":
		status = entities.JobStatusCompleted
	default:
		status = entities.JobStatusCompleted
	}
	
	conclusion := ""
	if resp.Conclusion != nil {
		conclusion = *resp.Conclusion
	}
	
	return &entities.Job{
		ID:           resp.ID,
		RunID:        resp.RunID,
		Name:         resp.Name,
		Status:       status,
		RunnerID:     resp.RunnerID,
		RunnerName:   resp.RunnerName,
		StartedAt:    resp.StartedAt,
		CompletedAt:  resp.CompletedAt,
		Conclusion:   conclusion,
		WorkflowName: workflowName,
		Repository:   repository,
	}
}