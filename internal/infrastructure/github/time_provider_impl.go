package github

import (
	"time"

	"github.com/VeyronSakai/gh-runner-monitor/internal/domain/repository"
)

// TimeProviderImpl implements the TimeProvider interface for production use
type TimeProviderImpl struct{}

// NewTimeProvider creates a new instance of TimeProviderImpl
func NewTimeProvider() repository.TimeProvider {
	return &TimeProviderImpl{}
}

// GetCurrentTime returns the current time
func (t *TimeProviderImpl) GetCurrentTime() time.Time {
	return time.Now()
}
