package presentation

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/usecase"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Column title constants
const (
	columnTitleRunnerName    = "Runner Name"
	columnTitleStatus        = "Status"
	columnTitleJobName       = "Job Name"
	columnTitleExecutionTime = "Execution Time"
)

// Column width constants
const (
	// Minimum widths for each column
	minWidthRunnerName = 15
	minWidthStatus     = 12
	minWidthJobName    = 20
	minWidthExecTime   = 15

	// Space reserved for borders and padding
	borderPadding = 10

	// Space reserved for header and footer in height calculation
	headerFooterHeight = 5

	// Proportions for distributing extra width
	ratioRunnerName = 0.25
	ratioStatus     = 0.15
	ratioJobName    = 0.45
	ratioExecTime   = 0.15

	// Default terminal size (fallback if WindowSizeMsg is not received)
	defaultTerminalWidth  = 80
	defaultTerminalHeight = 24
)

// Model represents the TUI application state
type Model struct {
	table          table.Model
	runnerMonitor  *usecase.RunnerMonitor
	owner          string
	repo           string
	org            string
	runners        []*entity.Runner
	jobs           []*entity.Job
	currentTime    time.Time
	lastUpdate     time.Time
	updateInterval time.Duration
	quitting       bool
	width          int
	height         int
	err            error
}

// NewModel creates a new TUI model
func NewModel(useCase *usecase.RunnerMonitor, owner, repo, org string, intervalSeconds int) *Model {
	// Start with minimum column widths - will be updated when WindowSizeMsg is received
	columns := []table.Column{
		{Title: columnTitleRunnerName, Width: minWidthRunnerName},
		{Title: columnTitleStatus, Width: minWidthStatus},
		{Title: columnTitleJobName, Width: minWidthJobName},
		{Title: columnTitleExecutionTime, Width: minWidthExecTime},
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

	// Calculate initial table height based on default terminal height
	tableHeight := getCalculatedTableHeight(defaultTerminalHeight)

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
		table.WithStyles(s),
	)

	return &Model{
		table:          t,
		runnerMonitor:  useCase,
		owner:          owner,
		repo:           repo,
		org:            org,
		updateInterval: time.Duration(intervalSeconds) * time.Second,
		width:          defaultTerminalWidth,
		height:         defaultTerminalHeight,
	}
}

// getCalculatedColumnWidths calculates column widths based on available terminal width
func getCalculatedColumnWidths(terminalWidth int) []table.Column {
	availableWidth := terminalWidth - borderPadding
	totalMinWidth := minWidthRunnerName + minWidthStatus + minWidthJobName + minWidthExecTime

	if availableWidth < totalMinWidth {
		// Terminal is too small, use minimum widths
		return []table.Column{
			{Title: columnTitleRunnerName, Width: minWidthRunnerName},
			{Title: columnTitleStatus, Width: minWidthStatus},
			{Title: columnTitleJobName, Width: minWidthJobName},
			{Title: columnTitleExecutionTime, Width: minWidthExecTime},
		}
	}

	remainingWidth := availableWidth - totalMinWidth

	// Distribute remaining width proportionally
	runnerNameExtra := int(float64(remainingWidth) * ratioRunnerName)
	statusExtra := int(float64(remainingWidth) * ratioStatus)
	jobNameExtra := int(float64(remainingWidth) * ratioJobName)
	execTimeExtra := remainingWidth - runnerNameExtra - statusExtra - jobNameExtra

	return []table.Column{
		{Title: columnTitleRunnerName, Width: minWidthRunnerName + runnerNameExtra},
		{Title: columnTitleStatus, Width: minWidthStatus + statusExtra},
		{Title: columnTitleJobName, Width: minWidthJobName + jobNameExtra},
		{Title: columnTitleExecutionTime, Width: minWidthExecTime + execTimeExtra},
	}
}

// getCalculatedTableHeight calculates table height based on terminal height
func getCalculatedTableHeight(terminalHeight int) int {
	height := terminalHeight - headerFooterHeight
	if height < 5 {
		return 5 // Minimum height to ensure table is usable
	}
	return height
}

// Init initializes the model and returns the initial command
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchData(),
		tea.Tick(m.updateInterval, func(t time.Time) tea.Msg {
			return t
		}),
	)
}
