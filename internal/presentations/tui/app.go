package tui

import (
	"fmt"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	useCase *use_cases.MonitorRunnersUseCase
	params  ports.MonitorParams
}

func NewApp(useCase *use_cases.MonitorRunnersUseCase, params ports.MonitorParams) *App {
	return &App{
		useCase: useCase,
		params:  params,
	}
}

func (a *App) Run() error {
	model := NewModel(a.useCase, a.params)
	ApplyTableStyles(&model.table)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}
	
	return nil
}