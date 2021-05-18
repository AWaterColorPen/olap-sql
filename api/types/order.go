package types

import "github.com/awatercolorpen/olap-sql/api/proto"

type OrderDirectionType string

func (o OrderDirectionType) ToEnum() proto.ORDER_DIRECTION_TYPE {
	if v, ok := proto.ORDER_DIRECTION_TYPE_value[string(o)]; ok {
		return proto.ORDER_DIRECTION_TYPE(v)
	}
	return proto.ORDER_DIRECTION_TYPE_ORDER_DIRECTION_UNKNOWN
}

func EnumToOrderDirectionType(o proto.ORDER_DIRECTION_TYPE) OrderDirectionType {
	return OrderDirectionType(o.String())
}

const (
	OrderDirectionTypeUnknown    OrderDirectionType = "ORDER_DIRECTION_UNKNOWN"
	OrderDirectionTypeAscending  OrderDirectionType = "ORDER_DIRECTION_ASCENDING"
	OrderDirectionTypeDescending OrderDirectionType = "ORDER_DIRECTION_DESCENDING"
)

type OrderBy struct {
	Name      string             `json:"name"`
	Direction OrderDirectionType `json:"direction"`
}

func (o *OrderBy) ToProto() *proto.OrderBy {
	return &proto.OrderBy{
		Name:      o.Name,
		Direction: o.Direction.ToEnum(),
	}
}

func ProtoToOrderBy(o *proto.OrderBy) *OrderBy {
	return &OrderBy{
		Name:      o.Name,
		Direction: EnumToOrderDirectionType(o.Direction),
	}
}
