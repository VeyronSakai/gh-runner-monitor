package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// GetRunners fetches runners for a repository or organization
func (c *Client) GetRunners(ctx context.Context, owner, repo, org string) ([]*entity.Runner, error) {
	var runners *runnersResponse
	var err error

	if org != "" {
		path := fmt.Sprintf("orgs/%s/actions/runners", org)
		runners, err = c.fetchRunners(path)
	} else {
		path := fmt.Sprintf("repos/%s/%s/actions/runners", owner, repo)
		runners, err = c.fetchRunners(path)
	}

	if err != nil {
		return nil, err
	}

	return c.convertRunners(runners), nil
}

func (c *Client) fetchRunners(path string) (*runnersResponse, error) {
	response, err := c.restClient.Request(http.MethodGet, path, nil)
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

func (c *Client) convertRunners(runners *runnersResponse) []*entity.Runner {
	result := make([]*entity.Runner, 0, len(runners.Runners))
	for _, r := range runners.Runners {
		status := entity.StatusOffline
		if r.Status == "online" {
			if r.Busy {
				status = entity.StatusActive
			} else {
				status = entity.StatusIdle
			}
		}

		labels := make([]string, 0, len(r.Labels))
		for _, l := range r.Labels {
			labels = append(labels, l.Name)
		}

		result = append(result, &entity.Runner{
			ID:        r.ID,
			Name:      r.Name,
			Status:    status,
			Labels:    labels,
			OS:        r.OS,
			UpdatedAt: time.Now(),
		})
	}

	return result
}
