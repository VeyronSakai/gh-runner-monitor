package ports

import (
	"context"
)

type MonitorRunnersInputPort interface {
	Execute(ctx context.Context, params MonitorParams) (*MonitorOutput, error)
}

type MonitorParams struct {
	Owner    string
	Repo     string
	Org      string
	IsOrgLevel bool
}

func NewMonitorParams(owner, repo string) MonitorParams {
	return MonitorParams{
		Owner:      owner,
		Repo:       repo,
		IsOrgLevel: false,
	}
}

func NewOrgMonitorParams(org string) MonitorParams {
	return MonitorParams{
		Org:        org,
		IsOrgLevel: true,
	}
}