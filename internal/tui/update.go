package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming events and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(msg.Height - 10)
		return m, nil
		
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, m.fetchData()
		}
		
	case tickMsg:
		return m, tea.Batch(
			m.fetchData(),
			tickCmd(m.updateInterval),
		)
		
	case dataMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.runners = msg.runners
			m.jobs = msg.jobs
			m.lastUpdate = time.Now()
			m.err = nil
			m.updateTableRows()
		}
		return m, nil
	}
	
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// updateTableRows updates the table with the current runner and job data
func (m *Model) updateTableRows() {
	rows := make([]table.Row, 0, len(m.runners))
	for _, runner := range m.runners {
		statusIcon := getStatusIcon(runner.Status)
		status := fmt.Sprintf("%s %s", statusIcon, runner.Status)
		
		jobName := "-"
		execTime := "-"
		
		// Find active job for this runner
		for _, job := range m.jobs {
			if job.RunnerID != nil && *job.RunnerID == runner.ID {
				jobName = fmt.Sprintf("%s (%s)", job.Name, job.WorkflowName)
				if job.StartedAt != nil {
					duration := time.Since(*job.StartedAt)
					execTime = formatDuration(duration)
				}
				break
			}
		}
		
		rows = append(rows, table.Row{
			runner.Name,
			status,
			jobName,
			execTime,
		})
	}
	m.table.SetRows(rows)
}

// fetchData fetches runners and jobs data from GitHub
func (m *Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		
		runners, err := m.client.GetRunners(ctx, m.owner, m.repo, m.org)
		if err != nil {
			return dataMsg{err: err}
		}
		
		jobs, err := m.client.GetActiveJobs(ctx, m.owner, m.repo, m.org)
		if err != nil {
			return dataMsg{err: err}
		}
		
		// Update runner status based on active jobs
		for _, runner := range runners {
			for _, job := range jobs {
				if job.RunnerID != nil && *job.RunnerID == runner.ID {
					runner.Status = models.StatusActive
					break
				}
			}
		}
		
		return dataMsg{runners: runners, jobs: jobs}
	}
}