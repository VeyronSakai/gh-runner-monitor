package presentation

import (
	"fmt"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// View returns the string representation of the model
func (m *Model) View() string {
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

// getStatusIcon returns the appropriate icon for the runner status
func getStatusIcon(status entity.RunnerStatus) string {
	switch status {
	case entity.StatusIdle:
		return "🟢"
	case entity.StatusActive:
		return "🟠"
	case entity.StatusOffline:
		return "⚫"
	default:
		return "❓"
	}
}

// formatDuration formats a time duration into a readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
