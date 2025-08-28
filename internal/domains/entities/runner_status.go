package entities

type RunnerStatus string

const (
	StatusIdle    RunnerStatus = "Idle"
	StatusActive  RunnerStatus = "Active"
	StatusOffline RunnerStatus = "Offline"
)

func (s RunnerStatus) String() string {
	return string(s)
}

func (s RunnerStatus) IsActive() bool {
	return s == StatusActive
}

func (s RunnerStatus) IsOffline() bool {
	return s == StatusOffline
}

func (s RunnerStatus) IsIdle() bool {
	return s == StatusIdle
}