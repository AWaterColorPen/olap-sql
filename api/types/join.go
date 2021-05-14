package types

import "github.com/awatercolorpen/olap-sql/api/proto"

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

func (j *Join) ToProto() *proto.Join {
	join := &proto.Join{
		Database1: j.Database1,
		Database2: j.Database2,
		Table1:    j.Table1,
		Table2:    j.Table2,
	}
	for _, v := range j.Filters {
		join.Filters = append(join.Filters, v.ToProto())
	}
	return join
}

func ProtoToJoin(j *proto.Join) *Join {
	join := &Join{
		Database1: j.Database1,
		Database2: j.Database2,
		Table1:    j.Table1,
		Table2:    j.Table2,
	}
	for _, v := range j.Filters {
		join.Filters = append(join.Filters, ProtoToFilter(v))
	}
	return join
}
