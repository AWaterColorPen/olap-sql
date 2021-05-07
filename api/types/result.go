package types

import "fmt"

type TableResult struct {
	Header []string        `json:"header"`
	Rows   [][]interface{} `json:"rows"`
}

func (t *TableResult) SetHeader(query *Query) {
	t.Header = append([]string{}, query.Dimensions...)
	t.Header = append(t.Header, query.Metrics...)
}

func (t *TableResult) AddRow( in map[string]interface{}) error {
	row := make([]interface{}, len(t.Header))
	for i, u := range t.Header {
		w, ok := in[u]
		if !ok {
			return fmt.Errorf("found no column %v", u)
		}
		row[i] = w
	}
	t.Rows = append(t.Rows, row)
	return nil
}

type SeriesResult struct {
	Header []string        `json:"header"`
	Rows   [][]interface{} `json:"rows"`
}