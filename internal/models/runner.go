package models

import "time"

type RunnerStatus string

const (
	StatusIdle    RunnerStatus = "Idle"
	StatusActive  RunnerStatus = "Active"
	StatusOffline RunnerStatus = "Offline"
)

type Runner struct {
	ID        int64
	Name      string
	Status    RunnerStatus
	Labels    []string
	OS        string
	UpdatedAt time.Time
}

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

type MonitorData struct {
	Runners []*Runner
	Jobs    []*Job
}