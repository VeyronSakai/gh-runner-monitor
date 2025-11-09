package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
)

// Data represents the structure of debug JSON data
type Data struct {
	CurrentTime time.Time        `json:"CurrentTime"`
	Runners     []*entity.Runner `json:"runners"`
	Jobs        []*entity.Job    `json:"jobs"`
}

// RunnerRepositoryImpl is a repository implementation that loads data from a JSON file
type RunnerRepositoryImpl struct {
	data *Data
}

// NewRunnerRepository creates a new debug repository from a JSON file
func NewRunnerRepository(jsonPath string) (repository.RunnerRepository, error) {
	data, err := loadData(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load debug data: %w", err)
	}

	return &RunnerRepositoryImpl{
		data: data,
	}, nil
}

// loadData loads debug data from a JSON file
func loadData(jsonPath string) (*Data, error) {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data Data
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &data, nil
}

// GetRunners returns the runners from the debug data
func (d *RunnerRepositoryImpl) FetchRunners(_ context.Context, _, _, _ string) ([]*entity.Runner, error) {
	return d.data.Runners, nil
}

// GetActiveJobs returns the jobs from the debug data
func (d *RunnerRepositoryImpl) FetchActiveJobs(_ context.Context, _, _, _ string) ([]*entity.Job, error) {
	return d.data.Jobs, nil
}

// GetCurrentTime returns the current time from the debug data
// This allows time to be mocked in debug mode
func (d *RunnerRepositoryImpl) GetCurrentTime() time.Time {
	return d.data.CurrentTime
}
