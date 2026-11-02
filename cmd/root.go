package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/VeyronSakai/gh-runner-monitor/internal/infrastructure/debug"
	"github.com/VeyronSakai/gh-runner-monitor/internal/infrastructure/github"
	"github.com/VeyronSakai/gh-runner-monitor/internal/presentation"
	"github.com/VeyronSakai/gh-runner-monitor/internal/usecase"
	tea "github.com/charmbracelet/bubbletea"
	ghrepo "github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

var (
	org       string
	repo      string
	interval  int
	debugPath string
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
	rootCmd.Flags().StringVar(&debugPath, "debug", "", "Debug mode: path to JSON file with mock runner data")
}

func runMonitor(_ *cobra.Command, _ []string) error {
	var runnerRepo repository.RunnerRepository
	var jobRepo repository.JobRepository
	var timeProvider repository.TimeProvider
	var err error

	// Check if debug mode is enabled
	if debugPath != "" {
		// Use debug repositories with JSON data
		data, err := debug.LoadDebugData(debugPath)
		if err != nil {
			return fmt.Errorf("failed to load debug data: %w", err)
		}
		runnerRepo = debug.NewRunnerRepository(data)
		jobRepo = debug.NewJobRepository(data)
		timeProvider = debug.NewTimeProvider(data)
	} else {
		// Create infrastructure layer (GitHub client)
		runnerRepo, err = github.NewRunnerRepository()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		jobRepo, err = github.NewJobRepository()
		if err != nil {
			return fmt.Errorf("failed to create GitHub job client: %w", err)
		}
		timeProvider = github.NewTimeProvider()
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
		// In debug mode, we don't need to fetch current repository
		if debugPath != "" {
			owner = "owner"
			repoName = "repo"
		} else {
			currentRepo, err := ghrepo.Current()
			if err != nil {
				return fmt.Errorf("not in a git repository and no --repo or --org flag specified")
			}
			owner = currentRepo.Owner
			repoName = currentRepo.Name
		}
	}

	// Create use case with dependencies
	monitorUseCase := usecase.NewRunnerMonitor(runnerRepo, jobRepo, timeProvider)

	// Create presentation layer (TUI) with use case
	model := presentation.NewModel(monitorUseCase, owner, repoName, orgName, interval)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
