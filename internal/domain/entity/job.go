package entity

import "time"

// Job represents a GitHub Actions workflow job
type Job struct {
	ID           int64
	RunID        int64
	Name         string
	Status       string
	RunnerID     *int64
	RunnerName   *string
	StartedAt    *time.Time
	WorkflowName string
	Repository   string
}

// IsRunning returns true if the job is currently running
func (j *Job) IsRunning() bool {
	return j.Status == "in_progress"
}

// IsQueued returns true if the job is queued
func (j *Job) IsQueued() bool {
	return j.Status == "queued"
}

// IsAssignedToRunner returns true if the job is assigned to a specific runner
func (j *Job) IsAssignedToRunner(runnerID int64) bool {
	return j.RunnerID != nil && *j.RunnerID == runnerID
}

// GetExecutionDuration returns the duration since the job started
func (j *Job) GetExecutionDuration() time.Duration {
	if j.StartedAt == nil {
		return 0
	}
	return time.Since(*j.StartedAt)
}
