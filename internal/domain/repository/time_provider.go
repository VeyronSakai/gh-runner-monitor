package repository

import "time"

// TimeProvider defines the interface for getting current time
type TimeProvider interface {
	// GetCurrentTime returns the current time (for mocking in tests/debug mode)
	GetCurrentTime() time.Time
}
