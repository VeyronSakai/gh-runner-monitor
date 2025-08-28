package entities

import "time"

type Runner struct {
	ID        int64
	Name      string
	Status    RunnerStatus
	Labels    []string
	GroupName string
	GroupID   int64
	OS        string
	UpdatedAt time.Time
}

func NewRunner(id int64, name string, status RunnerStatus) *Runner {
	return &Runner{
		ID:        id,
		Name:      name,
		Status:    status,
		Labels:    []string{},
		UpdatedAt: time.Now(),
	}
}

func (r *Runner) UpdateStatus(status RunnerStatus) {
	r.Status = status
	r.UpdatedAt = time.Now()
}

func (r *Runner) HasLabel(label string) bool {
	for _, l := range r.Labels {
		if l == label {
			return true
		}
	}
	return false
}