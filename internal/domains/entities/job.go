package entities

import (
	"time"
)

type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusCompleted  JobStatus = "completed"
)

type Job struct {
	ID          int64
	RunID       int64
	Name        string
	Status      JobStatus
	RunnerID    *int64
	RunnerName  *string
	StartedAt   *time.Time
	CompletedAt *time.Time
	Conclusion  string
	WorkflowName string
	Repository   string
}

func NewJob(id int64, runID int64, name string, status JobStatus) *Job {
	return &Job{
		ID:     id,
		RunID:  runID,
		Name:   name,
		Status: status,
	}
}

func (j *Job) IsRunning() bool {
	return j.Status == JobStatusInProgress
}

func (j *Job) GetExecutionTime() time.Duration {
	if j.StartedAt == nil {
		return 0
	}
	
	endTime := time.Now()
	if j.CompletedAt != nil {
		endTime = *j.CompletedAt
	}
	
	return endTime.Sub(*j.StartedAt)
}

func (j *Job) GetRunnerID() int64 {
	if j.RunnerID == nil {
		return 0
	}
	return *j.RunnerID
}