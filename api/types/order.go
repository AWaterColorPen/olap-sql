package types

import (
	"fmt"
)

type OrderDirectionType string

const (
	OrderDirectionTypeUnknown    OrderDirectionType = "ORDER_DIRECTION_UNKNOWN"
	OrderDirectionTypeAscending  OrderDirectionType = "ORDER_DIRECTION_ASCENDING"
	OrderDirectionTypeDescending OrderDirectionType = "ORDER_DIRECTION_DESCENDING"
)

type OrderBy struct {
	Table     string             `json:"table"`
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
