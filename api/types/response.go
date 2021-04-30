package types

import "github.com/awatercolorpen/olap-sql/api/proto"

type Response struct {
	Rows []*Item `json:"rows"`
}

func (r *Response) ToProto() *proto.Response {
	response := &proto.Response{}
	for _, v := range r.Rows {
		response.Rows = append(response.Rows, v.ToProto())
	}
	return response
}

func ProtoToResponse(response *proto.Response) *Response {
	r := &Response{}
	for _, v := range response.Rows {
		r.Rows = append(r.Rows, ProtoToItem(v))
	}
	return r
}

type Item struct {
	Values map[string]string `json:"values"`
}

func (item *Item) ToProto() *proto.Item {
	return &proto.Item{
		Values: item.Values,
	}
}

func ProtoToItem(item *proto.Item) *Item {
	return &Item{item.Values}
}
