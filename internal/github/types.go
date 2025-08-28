package github

import "time"

// API response types
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