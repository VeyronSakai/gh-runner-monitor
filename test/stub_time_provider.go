package test

import "time"

// StubTimeProvider is a stub implementation of repository.TimeProvider for testing.
type StubTimeProvider struct {
	// CurrentTime is the time that will be returned by GetCurrentTime
	CurrentTime time.Time
}

// GetCurrentTime returns the configured current time for testing
func (s *StubTimeProvider) GetCurrentTime() time.Time {
	if s.CurrentTime.IsZero() {
		return time.Now()
	}
	return s.CurrentTime
}
