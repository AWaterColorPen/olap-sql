package types

type JoinOn struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

type Join struct {
	DataSource1 *DataSource       `json:"datasource1"`
	DataSource2 *DataSource       `json:"datasource2"`
	On          []*JoinOn         `json:"on"`
	Filters     []*Filter         `json:"filters"`
	TableSqlMap map[string]string `json:"tablesqlsmap"`
}

func (j *Join) Statement() (string, error) {
	return "", nil
}
