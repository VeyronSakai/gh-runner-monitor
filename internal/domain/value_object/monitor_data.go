package value_object

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/entity"
)

// MonitorData represents the data for the runner monitor
type MonitorData struct {
	CurrentTime time.Time
	Runners     []*entity.Runner
	Jobs        []*entity.Job
}
