package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	domainrepo "github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
	"github.com/cli/go-gh/v2/pkg/api"
)

// RunnerRepositoryImpl implements the RunnerRepository interface using GitHub API
type RunnerRepositoryImpl struct {
	restClient *api.RESTClient
}

// NewRunnerRepository creates a new instance of RunnerRepositoryImpl
func NewRunnerRepository() (domainrepo.RunnerRepository, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w\nPlease run 'gh auth login' to authenticate with GitHub", err)
	}

	return &RunnerRepositoryImpl{
		restClient: restClient,
	}, nil
}

// FetchRunners retrieves all runners for a repository or organization
func (r *RunnerRepositoryImpl) FetchRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error) {
	path := r.getRunnersPath(owner, repo, org)
	runners, err := r.requestGetRunners(path)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Runner, 0, len(runners.Runners))
	for _, runner := range runners.Runners {
		status := entity.StatusOffline
		if runner.Status == "online" {
			if runner.Busy {
				status = entity.StatusActive
			} else {
				status = entity.StatusIdle
			}
		}

		labels := make([]string, 0, len(runner.Labels))
		for _, l := range runner.Labels {
			labels = append(labels, l.Name)
		}

		result = append(result, &entity.Runner{
			ID:        runner.ID,
			Name:      runner.Name,
			Status:    status,
			Labels:    labels,
			OS:        runner.OS,
			UpdatedAt: time.Now(),
		})
	}
	return result, nil
}

// getRunnersPath constructs the API path for fetching runners
func (r *RunnerRepositoryImpl) getRunnersPath(owner, repo, org string) string {
	if org != "" {
		return fmt.Sprintf("orgs/%s/actions/runners", org)
	}
	return fmt.Sprintf("repos/%s/%s/actions/runners", owner, repo)
}

// requestGetRunners fetches runners from GitHub API
func (r *RunnerRepositoryImpl) requestGetRunners(path string) (*runnersResponse, error) {
	response, err := r.restClient.Request(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request runners: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var runners runnersResponse
	if err := json.NewDecoder(response.Body).Decode(&runners); err != nil {
		return nil, fmt.Errorf("failed to decode runners response: %w", err)
	}

	return &runners, nil
}
