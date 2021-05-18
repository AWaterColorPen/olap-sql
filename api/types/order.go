package types

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

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

func (o *OrderBy) Expression() (string, error) {
	switch o.Direction {
	case OrderDirectionTypeAscending:
		return fmt.Sprintf("%v ASC", o.Name), nil
	case OrderDirectionTypeDescending:
		return fmt.Sprintf("%v DESC", o.Name), nil
	default:
		return "", fmt.Errorf("not supported order direction type %v", o.Direction)
	}
}

func (o *OrderBy) Alias() (string, error) {
	return "", fmt.Errorf("order by is unsupported alias method")
}

func (o *OrderBy) Statement() (string, error) {
	return o.Expression()
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
