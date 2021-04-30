package types

import "github.com/awatercolorpen/olap-sql/api/proto"

type Response []Item

func (r Response) ToProto() *proto.Response {
	response := &proto.Response{}
	for _, v := range r {
		response.Rows = append(response.Rows, v.ToProto())
	}
	return response
}

func ProtoToResponse(response *proto.Response) Response {
	r := Response{}
	for _, v := range response.Rows {
		r = append(r, ProtoToItem(v))
	}
	return r
}

type Item map[string]string

func (item Item) ToProto() *proto.Item {
	return &proto.Item{
		Values: item,
	}
}

func ProtoToItem(item *proto.Item) Item {
	return item.Values
}
