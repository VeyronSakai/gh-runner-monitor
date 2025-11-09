package debug

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
)

// TimeProviderImpl is a time provider implementation that returns time from debug data
type TimeProviderImpl struct {
	data *Data
}

// NewTimeProvider creates a new debug time provider from loaded data
func NewTimeProvider(data *Data) repository.TimeProvider {
	return &TimeProviderImpl{
		data: data,
	}
}

// GetCurrentTime returns the current time from the debug data
// This allows time to be mocked in debug mode
func (t *TimeProviderImpl) GetCurrentTime() time.Time {
	return t.data.CurrentTime
}
