package types

import "github.com/awatercolorpen/olap-sql/api/proto"

type Join struct {
	Table1 string `json:"table1"`
	Key1   string `json:"key1"`
	Table2 string `json:"table2"`
	Key2   string `json:"key2"`
}

func (j *Join) ToProto() *proto.Join {
	return &proto.Join{
		Table1: j.Table1,
		Key1:   j.Key1,
		Table2: j.Table2,
		Key2:   j.Key2,
	}
}

func ProtoToJoin(j *proto.Join) *Join {
	return &Join{
		Table1: j.Table1,
		Key1:   j.Key1,
		Table2: j.Table2,
		Key2:   j.Key2,
	}
}
