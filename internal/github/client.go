package github

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
)

// Client represents a GitHub API client
type Client struct {
	restClient *api.RESTClient
}

// NewClient creates a new GitHub API client
func NewClient() (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}
	
	return &Client{
		restClient: restClient,
	}, nil
}