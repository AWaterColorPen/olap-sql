package types

type TimeInterval struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func (t *TimeInterval) ToFilter() *Filter {
	filter1 := &Filter{OperatorType: FilterOperatorTypeGreaterEquals, Name: t.Name, Value: []interface{}{t.Start}}
	filter2 := &Filter{OperatorType: FilterOperatorTypeLess, Name: t.Name, Value: []interface{}{t.End}}
	return &Filter{
		OperatorType: FilterOperatorTypeAnd,
		Children:     []*Filter{filter1, filter2},
	}
}

type Query struct {
	DataSetName  string        `json:"data_set_name"`
	TimeInterval *TimeInterval `json:"time_interval"`
	Metrics      []string      `json:"metrics"`
	Dimensions   []string      `json:"dimensions"`
	Filters      []*Filter     `json:"filters"`
	Orders       []*OrderBy    `json:"orders"`
	Limit        *Limit        `json:"limit"`
	Sql          string        `json:"Sql"`
}

func (q *Query) TranslateTimeIntervalToFilter() {
	if q.TimeInterval != nil && q.TimeInterval.Name != "" && q.TimeInterval.Start != "" && q.TimeInterval.End != "" {
		q.Filters = append(q.Filters, q.TimeInterval.ToFilter())
	}
}
