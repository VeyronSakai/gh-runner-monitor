package presentation

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/usecase"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Column title constants
const (
	columnTitleRunnerName    = "Runner"
	columnTitleStatus        = "Status"
	columnTitleLabels        = "Labels"
	columnTitleJobName       = "Job Name"
	columnTitleExecutionTime = "Time"
)

// Column width constants
const (
	// Minimum widths for each column
	minRunnerNameWidth = 10
	minLabelsWidth     = 10
	minJobNameWidth    = 10

	statusWidth   = 10
	execTimeWidth = 10

	// Space reserved for borders and padding
	borderPadding = 10

	// Space reserved for header and footer in height calculation
	headerFooterHeight = 5

	// Proportions for distributing extra width
	ratioRunnerName = 0.10
	ratioLabels     = 0.35
	ratioJobName    = 0.55

	// Default terminal size (fallback if WindowSizeMsg is not received)
	defaultTerminalWidth  = 80
	defaultTerminalHeight = 24
)

// Model represents the TUI application state
type Model struct {
	table          table.Model
	spinner        spinner.Model
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
	loading        bool
	width          int
	height         int
	err            error
}

// NewModel creates a new TUI model
func NewModel(useCase *usecase.RunnerMonitor, owner, repo, org string, intervalSeconds int) *Model {
	// Start with minimum column widths - will be updated when WindowSizeMsg is received
	columns := []table.Column{
		{Title: columnTitleRunnerName, Width: minRunnerNameWidth},
		{Title: columnTitleStatus, Width: statusWidth},
		{Title: columnTitleLabels, Width: minLabelsWidth},
		{Title: columnTitleJobName, Width: minJobNameWidth},
		{Title: columnTitleExecutionTime, Width: execTimeWidth},
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

	// Initialize spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Model{
		table:          t,
		spinner:        sp,
		runnerMonitor:  useCase,
		owner:          owner,
		repo:           repo,
		org:            org,
		updateInterval: time.Duration(intervalSeconds) * time.Second,
		loading:        true,
		width:          defaultTerminalWidth,
		height:         defaultTerminalHeight,
	}
}

// getCalculatedColumnWidths calculates column widths based on available terminal width
func getCalculatedColumnWidths(terminalWidth int) []table.Column {
	availableWidth := terminalWidth - borderPadding
	totalMinWidth := minRunnerNameWidth + statusWidth + minLabelsWidth + minJobNameWidth + execTimeWidth

	if availableWidth < totalMinWidth {
		// Terminal is too small, use minimum widths
		return []table.Column{
			{Title: columnTitleRunnerName, Width: minRunnerNameWidth},
			{Title: columnTitleStatus, Width: statusWidth},
			{Title: columnTitleLabels, Width: minLabelsWidth},
			{Title: columnTitleJobName, Width: minJobNameWidth},
			{Title: columnTitleExecutionTime, Width: execTimeWidth},
		}
	}

	remainingWidth := availableWidth - totalMinWidth

	// Distribute remaining width proportionally
	runnerNameExtra := int(float64(remainingWidth) * ratioRunnerName)
	labelsExtra := int(float64(remainingWidth) * ratioLabels)
	jobNameExtra := int(float64(remainingWidth) * ratioJobName)

	return []table.Column{
		{Title: columnTitleRunnerName, Width: minRunnerNameWidth + runnerNameExtra},
		{Title: columnTitleStatus, Width: statusWidth},
		{Title: columnTitleLabels, Width: minLabelsWidth + labelsExtra},
		{Title: columnTitleJobName, Width: minJobNameWidth + jobNameExtra},
		{Title: columnTitleExecutionTime, Width: execTimeWidth},
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
		m.spinner.Tick,
		m.fetchData(),
		tea.Tick(m.updateInterval, func(t time.Time) tea.Msg {
			return t
		}),
	)
}
