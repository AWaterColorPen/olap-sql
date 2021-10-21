package types

type JoinOn struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

type Join struct {
	Database1 string    `json:"data_base1"`
	Database2 string    `json:"data_base2"`
	Table1    string    `json:"table1"`
	Table2    string    `json:"table2"`
	On        []*JoinOn `json:"on"`
	Filters   []*Filter `json:"filters"`
}
