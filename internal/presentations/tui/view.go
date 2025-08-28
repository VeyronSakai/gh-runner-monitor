package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)
	
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))
	
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)
	
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)
		
	headerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)
)

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	
	var sb strings.Builder
	
	title := "🏃 GitHub Actions Runner Monitor"
	sb.WriteString(titleStyle.Render(title))
	sb.WriteString("\n")
	
	target := "Repository"
	targetValue := fmt.Sprintf("%s/%s", m.params.Owner, m.params.Repo)
	if m.params.IsOrgLevel {
		target = "Organization"
		targetValue = m.params.Org
	}
	info := fmt.Sprintf("%s: %s | Last Update: %s", target, targetValue, m.lastUpdate.Format("15:04:05"))
	sb.WriteString(infoStyle.Render(info))
	sb.WriteString("\n\n")
	
	if m.err != nil {
		sb.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		sb.WriteString("\n\n")
	}
	
	sb.WriteString(m.table.View())
	sb.WriteString("\n")
	
	help := "Keys: [q]uit | [r]efresh | [↑/↓] navigate"
	sb.WriteString(helpStyle.Render(help))
	
	return sb.String()
}