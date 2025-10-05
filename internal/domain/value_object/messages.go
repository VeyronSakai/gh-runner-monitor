package value_object

// DataMsg contains the fetched monitoring data
type DataMsg struct {
	Data *MonitorData
	Err  error
}
