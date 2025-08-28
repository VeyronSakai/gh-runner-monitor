package tui

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent on a timer to refresh the data
type tickMsg time.Time

// dataMsg contains the fetched runners and jobs data
type dataMsg struct {
	runners []*models.Runner
	jobs    []*models.Job
	err     error
}

// tickCmd returns a command that sends a tick message after the interval
func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}