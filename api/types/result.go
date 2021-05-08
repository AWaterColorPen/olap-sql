package types

type Result struct {
	Dimensions []string                 `json:"dimensions"`
	Source     []map[string]interface{} `json:"source"`
}

func (r *Result) SetDimensions(query *Query) {
	r.Dimensions = append([]string{}, query.Dimensions...)
	r.Dimensions = append(r.Dimensions, query.Metrics...)
}

func (r *Result) AddSource(in map[string]interface{}) error {
	r.Source = append(r.Source, in)
	return nil
}
