package value_object

import (
	"time"
)

// TickMsg is sent on a timer to refresh the data
type TickMsg time.Time

// DataMsg contains the fetched monitoring data
type DataMsg struct {
	Data *MonitorData
	Err  error
}
