package types

import "github.com/awatercolorpen/olap-sql/api/proto"

type Limit struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

func (l *Limit) ToProto() *proto.Limit {
	return &proto.Limit{
		Limit: l.Limit,
		Offset: l.Offset,
	}
}

func ProtoToLimit(l *proto.Limit) *Limit {
	return &Limit{
		Limit: l.Limit,
		Offset: l.Offset,
	}
}
