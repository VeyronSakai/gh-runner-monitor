package presentation

import (
	"context"
	"fmt"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/value_object"
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

	case time.Time:
		return m, tea.Batch(
			m.fetchData(),
			tea.Tick(m.updateInterval, func(t time.Time) tea.Msg {
				return t
			}),
		)

	case value_object.DataMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.runners = msg.Data.Runners
			m.jobs = msg.Data.Jobs
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
			if job.IsAssignedToRunner(runner.ID) {
				jobName = fmt.Sprintf("%s (%s)", job.Name, job.WorkflowName)
				execTime = formatDuration(job.GetExecutionDuration())
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

// fetchData fetches runners and jobs data using the use case
func (m *Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		data, err := m.useCase.Execute(ctx, m.owner, m.repo, m.org)
		if err != nil {
			return value_object.DataMsg{Err: err}
		}

		return value_object.DataMsg{Data: data}
	}
}
