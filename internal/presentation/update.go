package presentation

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/value_object"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming events and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(getCalculatedTableHeight(msg.Height))
		m.updateColumnWidths()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, m.fetchData()
		case "enter", "return":
			return m, m.openJobLog()
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
			m.currentTime = msg.Data.CurrentTime
			m.lastUpdate = time.Now()
			m.err = nil
			m.updateTableRows()
		}
		return m, nil
	}

	var curCmd tea.Cmd
	m.table, curCmd = m.table.Update(msg)

	return m, curCmd
}

// updateTableRows updates the table with the current runner and job data
func (m *Model) updateTableRows() {
	rows := make([]table.Row, 0, len(m.runners))
	for _, runner := range m.runners {
		statusIcon := getStatusIcon(runner.Status)
		status := fmt.Sprintf("%s %s", statusIcon, runner.Status)

		jobName := "-"
		execTime := "-"

		// Find active job for this runner by matching runner name
		for _, job := range m.jobs {
			if job.IsAssignedToRunner(runner.Name) {
				jobName = fmt.Sprintf("%s (%s)", job.Name, job.WorkflowName)
				execTime = formatDuration(job.GetExecutionDurationAt(m.currentTime))
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

// openJobLog opens the job log page in the browser for the currently selected row
func (m *Model) openJobLog() tea.Cmd {
	return func() tea.Msg {
		// Get the selected row index
		selectedRow := m.table.Cursor()
		if selectedRow < 0 || selectedRow >= len(m.runners) {
			return nil
		}

		// Get the runner for the selected row
		runner := m.runners[selectedRow]

		// Find the active job for this runner by matching runner name
		var jobURL string
		for _, job := range m.jobs {
			if job.IsAssignedToRunner(runner.Name) {
				jobURL = job.HtmlUrl
				break
			}
		}

		// If no job URL found, do nothing
		if jobURL == "" {
			return nil
		}

		// Open the URL in the default browser
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", jobURL)
		case "linux":
			cmd = exec.Command("xdg-open", jobURL)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", jobURL)
		default:
			return nil
		}

		if err := cmd.Start(); err != nil {
			return value_object.DataMsg{Err: err}
		}

		return nil
	}
}

// fetchData fetches runners and jobs data using the use case
func (m *Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		data, err := m.runnerMonitor.Execute(ctx, m.owner, m.repo, m.org)
		if err != nil {
			return value_object.DataMsg{Err: err}
		}

		return value_object.DataMsg{Data: data}
	}
}

// updateColumnWidths adjusts column widths based on terminal width
func (m *Model) updateColumnWidths() {
	columns := getCalculatedColumnWidths(m.width)
	m.table.SetColumns(columns)
}
