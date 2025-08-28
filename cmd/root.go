package cmd

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/services"
	"github.com/VeyronSakai/gh-runner-monitor/internal/infrastructures/factories"
	"github.com/VeyronSakai/gh-runner-monitor/internal/presentations/tui"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases"
	"github.com/VeyronSakai/gh-runner-monitor/internal/use_cases/ports"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&org, "org", "", "Monitor runners for an organization")
	rootCmd.Flags().StringVar(&repo, "repo", "", "Monitor runners for a specific repository (owner/repo)")
	rootCmd.Flags().IntVar(&interval, "interval", 5, "Update interval in seconds")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	var params ports.MonitorParams
	
	if org != "" {
		params = ports.NewOrgMonitorParams(org)
	} else if repo != "" {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format. Use owner/repo")
		}
		params = ports.NewMonitorParams(parts[0], parts[1])
	} else {
		currentRepo, err := repository.Current()
		if err != nil {
			return fmt.Errorf("not in a git repository and no --repo or --org flag specified")
		}
		params = ports.NewMonitorParams(currentRepo.Owner, currentRepo.Name)
	}
	
	factory, err := factories.NewRepositoryFactory()
	if err != nil {
		return fmt.Errorf("failed to create repository factory: %w", err)
	}
	
	runnerRepo := factory.CreateRunnerRepository()
	jobRepo := factory.CreateJobRepository()
	monitorService := services.NewRunnerMonitorService()
	
	useCase := use_cases.NewMonitorRunnersUseCase(runnerRepo, jobRepo, monitorService)
	
	app := tui.NewApp(useCase, params)
	
	if err := app.Run(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}
	
	return nil
}