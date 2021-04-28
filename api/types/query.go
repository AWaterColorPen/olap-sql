package types

import "time"

type Query struct {
	Metrics      []string      `json:"metrics"`
	Dimensions   []string      `json:"dimensions"`
	Filters      []*Filter     `json:"filters"`
	TimeInterval *TimeInterval `json:"time_interval"`
	DataSet      string        `json:"data_set"`
}

type TimeInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
