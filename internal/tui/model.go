package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/github"
	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table          table.Model
	client         *github.Client
	owner          string
	repo           string
	org            string
	runners        []*models.Runner
	jobs           []*models.Job
	lastUpdate     time.Time
	updateInterval time.Duration
	err            error
	quitting       bool
	width          int
	height         int
}

func NewModel(client *github.Client, owner, repo, org string) Model {
	columns := []table.Column{
		{Title: "Runner Name", Width: 25},
		{Title: "Status", Width: 12},
		{Title: "Job Name", Width: 35},
		{Title: "Execution Time", Width: 15},
	}
	
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithStyles(s),
	)
	
	return Model{
		table:          t,
		client:         client,
		owner:          owner,
		repo:           repo,
		org:            org,
		updateInterval: 5 * time.Second,
	}
}

type tickMsg time.Time

type dataMsg struct {
	runners []*models.Runner
	jobs    []*models.Job
	err     error
}

func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) fetchData() tea.Cmd {
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

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchData(),
		tickCmd(m.updateInterval),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	
	var header string
	if m.org != "" {
		header = fmt.Sprintf("GitHub Runners Monitor - Organization: %s\n", m.org)
	} else {
		header = fmt.Sprintf("GitHub Runners Monitor - Repository: %s/%s\n", m.owner, m.repo)
	}
	header += fmt.Sprintf("Last Updated: %s | Press 'q' to quit, 'r' to refresh\n\n", 
		m.lastUpdate.Format("15:04:05"))
	
	if m.err != nil {
		return header + fmt.Sprintf("\nError: %v\n", m.err)
	}
	
	return header + m.table.View()
}

func getStatusIcon(status models.RunnerStatus) string {
	switch status {
	case models.StatusIdle:
		return "🟢"
	case models.StatusActive:
		return "🟠"
	case models.StatusOffline:
		return "⚫"
	default:
		return "❓"
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}