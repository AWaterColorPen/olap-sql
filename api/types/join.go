package types

type JoinOn struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

type Join struct {
	DataSource1 *DataSource `json:"datasource1"`
	DataSource2 *DataSource `json:"datasource2"`
	JoinType    string      `json:"join_type"`
	On          []*JoinOn   `json:"on"`
	Filters     []*Filter   `json:"filters"`
}

func (j *Join) GetJoinType() string {
	if len(j.JoinType) != 0 {
		return j.JoinType
	}
	return "LEFT JOIN"
}
