package tui

import (
	"context"
	"fmt"
	"time"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	table          table.Model
	useCase        *use_cases.MonitorRunnersUseCase
	params         ports.MonitorParams
	runners        []*ports.RunnerOutput
	lastUpdate     time.Time
	updateInterval time.Duration
	err            error
	quitting       bool
	width          int
	height         int
}

func NewModel(useCase *use_cases.MonitorRunnersUseCase, params ports.MonitorParams) Model {
	columns := []table.Column{
		{Title: "Runner Name", Width: 25},
		{Title: "Status", Width: 12},
		{Title: "Job Name", Width: 35},
		{Title: "Execution Time", Width: 15},
	}
	
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	
	return Model{
		table:          t,
		useCase:        useCase,
		params:         params,
		updateInterval: 5 * time.Second,
	}
}

type tickMsg time.Time

type runnersMsg struct {
	output *ports.MonitorOutput
	err    error
}

func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) fetchRunners() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		output, err := m.useCase.Execute(ctx, m.params)
		return runnersMsg{output: output, err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchRunners(),
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
			return m, m.fetchRunners()
		}
		
	case tickMsg:
		return m, tea.Batch(
			m.fetchRunners(),
			tickCmd(m.updateInterval),
		)
		
	case runnersMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.runners = msg.output.Runners
			m.lastUpdate = msg.output.UpdatedAt
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
		status := fmt.Sprintf("%s %s", runner.StatusIcon, runner.Status)
		jobName := runner.JobName
		if jobName == "" {
			jobName = "-"
		}
		execTime := runner.ExecutionTime
		if execTime == "" {
			execTime = "-"
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