package types

type Result struct {
	Header []string        `json:"header"`
	Rows   [][]interface{} `json:"rows"`
}
