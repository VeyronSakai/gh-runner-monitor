package ports

import (
	"fmt"
	"time"
)

type MonitorOutput struct {
	Runners   []*RunnerOutput
	UpdatedAt time.Time
	ErrorMsg  string
}

type RunnerOutput struct {
	ID            int64
	Name          string
	Status        string
	StatusColor   string
	StatusIcon    string
	JobName       string
	ExecutionTime string
	Labels        []string
	OS            string
}

func NewRunnerOutput(id int64, name, status string) *RunnerOutput {
	output := &RunnerOutput{
		ID:     id,
		Name:   name,
		Status: status,
	}
	
	switch status {
	case "Idle":
		output.StatusColor = "green"
		output.StatusIcon = "🟢"
	case "Active":
		output.StatusColor = "orange"
		output.StatusIcon = "🟠"
	case "Offline":
		output.StatusColor = "gray"
		output.StatusIcon = "⚫"
	default:
		output.StatusColor = "white"
		output.StatusIcon = "⚪"
	}
	
	return output
}

func (r *RunnerOutput) SetJob(jobName string, executionTime time.Duration) {
	r.JobName = jobName
	if executionTime > 0 {
		r.ExecutionTime = formatDuration(executionTime)
	}
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}