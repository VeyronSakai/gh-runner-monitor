package cmd

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/VeyronSakai/gh-runner-monitor/internal/github"
	"github.com/VeyronSakai/gh-runner-monitor/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

var (
	org      string
	repo     string
	interval int
)

var rootCmd = &cobra.Command{
	Use:   "gh-runner-monitor",
	Short: "Monitor GitHub Actions self-hosted runners in real-time",
	Long: `GitHub Actions Runner Monitor is a TUI tool that displays the status 
of self-hosted runners with their current jobs and execution times.`,
	RunE: runMonitor,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Explicitly ignore the error from Fprintln as we're already exiting
		// and there's nothing meaningful we can do if stderr write fails
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&org, "org", "", "Monitor runners for an organization")
	rootCmd.Flags().StringVar(&repo, "repo", "", "Monitor runners for a specific repository (owner/repo)")
	rootCmd.Flags().IntVar(&interval, "interval", 5, "Update interval in seconds")
}

func runMonitor(_ *cobra.Command, _ []string) error {
	client, err := github.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}
	
	var owner, repoName, orgName string
	
	if org != "" {
		orgName = org
	} else if repo != "" {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format. Use owner/repo")
		}
		owner = parts[0]
		repoName = parts[1]
	} else {
		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("not in a git repository and no --repo or --org flag specified")
		}
		owner = currentRepo.Owner
		repoName = currentRepo.Name
	}
	
	model := tui.NewModel(client, owner, repoName, orgName)
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	
	return nil
}