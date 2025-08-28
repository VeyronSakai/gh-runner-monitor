package tui

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/github"
	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI application state
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

// NewModel creates a new TUI model
func NewModel(client *github.Client, owner, repo, org string) *Model {
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
	
	return &Model{
		table:          t,
		client:         client,
		owner:          owner,
		repo:           repo,
		org:            org,
		updateInterval: 5 * time.Second,
	}
}

// Init initializes the model and returns the initial command
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchData(),
		tickCmd(m.updateInterval),
	)
}