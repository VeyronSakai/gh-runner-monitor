package entity

import "time"

// RunnerStatus represents the status of a runner
type RunnerStatus string

const (
	StatusIdle    RunnerStatus = "Idle"
	StatusActive  RunnerStatus = "Active"
	StatusOffline RunnerStatus = "Offline"
)

// Runner represents a GitHub Actions self-hosted runner
type Runner struct {
	ID        int64
	Name      string
	Status    RunnerStatus
	Labels    []string
	OS        string
	UpdatedAt time.Time
}

// IsOnline returns true if the runner is online (idle or active)
func (r *Runner) IsOnline() bool {
	return r.Status == StatusIdle || r.Status == StatusActive
}

// IsActive returns true if the runner is active
func (r *Runner) IsActive() bool {
	return r.Status == StatusActive
}
