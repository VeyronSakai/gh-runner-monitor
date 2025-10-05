package value_object

import (
	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// MonitorData represents the monitoring result
type MonitorData struct {
	Runners []*entity.Runner
	Jobs    []*entity.Job
}
