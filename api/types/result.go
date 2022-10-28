package types

type Result struct {
	Dimensions []string         `json:"dimensions"`
	Source     []map[string]any `json:"source"`
}

func (r *Result) SetDimensions(query *Query) {
	r.Dimensions = append([]string{}, query.Dimensions...)
	r.Dimensions = append(r.Dimensions, query.Metrics...)
}

func (r *Result) AddSource(in map[string]any) error {
	r.Source = append(r.Source, in)
	return nil
}
