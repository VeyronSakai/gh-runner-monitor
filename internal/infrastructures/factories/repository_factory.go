package factories

import (
	"fmt"
	"github.com/VeyronSakai/gh-runner-monitor/internal/domains/repositories"
	"github.com/VeyronSakai/gh-runner-monitor/internal/infrastructures/repositories/github"
)

type RepositoryFactory struct {
	githubClient *github.Client
}

func NewRepositoryFactory() (*RepositoryFactory, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}
	
	return &RepositoryFactory{
		githubClient: client,
	}, nil
}

func (f *RepositoryFactory) CreateRunnerRepository() repositories.RunnerRepository {
	return github.NewRunnerRepository(f.githubClient)
}

func (f *RepositoryFactory) CreateJobRepository() repositories.JobRepository {
	return github.NewJobRepository(f.githubClient)
}