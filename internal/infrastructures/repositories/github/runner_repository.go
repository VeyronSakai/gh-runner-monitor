package github

import (
	"context"
	"fmt"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/entities"
)

type RunnerRepository struct {
	client *Client
}

func NewRunnerRepository(client *Client) *RunnerRepository {
	return &RunnerRepository{
		client: client,
	}
}

func (r *RunnerRepository) ListRunners(ctx context.Context, owner, repo string) ([]*entities.Runner, error) {
	response, err := r.client.GetRepoRunners(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo runners: %w", err)
	}
	
	runners := make([]*entities.Runner, 0, len(response.Runners))
	for _, runnerResp := range response.Runners {
		runner := r.convertToEntity(runnerResp)
		runners = append(runners, runner)
	}
	
	return runners, nil
}

func (r *RunnerRepository) ListOrgRunners(ctx context.Context, org string) ([]*entities.Runner, error) {
	response, err := r.client.GetOrgRunners(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to get org runners: %w", err)
	}
	
	runners := make([]*entities.Runner, 0, len(response.Runners))
	for _, runnerResp := range response.Runners {
		runner := r.convertToEntity(runnerResp)
		runners = append(runners, runner)
	}
	
	return runners, nil
}

func (r *RunnerRepository) GetRunner(ctx context.Context, owner, repo string, runnerID int64) (*entities.Runner, error) {
	runners, err := r.ListRunners(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	
	for _, runner := range runners {
		if runner.ID == runnerID {
			return runner, nil
		}
	}
	
	return nil, fmt.Errorf("runner with ID %d not found", runnerID)
}

func (r *RunnerRepository) convertToEntity(resp RunnerResponse) *entities.Runner {
	var status entities.RunnerStatus
	
	if resp.Status == "offline" {
		status = entities.StatusOffline
	} else if resp.Busy {
		status = entities.StatusActive
	} else {
		status = entities.StatusIdle
	}
	
	labels := make([]string, 0, len(resp.Labels))
	for _, label := range resp.Labels {
		labels = append(labels, label.Name)
	}
	
	return &entities.Runner{
		ID:     resp.ID,
		Name:   resp.Name,
		Status: status,
		Labels: labels,
		OS:     resp.OS,
	}
}