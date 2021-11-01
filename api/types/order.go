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
	Table         string             `json:"table"`
	Name          string             `json:"name"`
	FieldProperty FieldProperty      `json:"field_property"`
	Direction     OrderDirectionType `json:"direction"`
}

func (o *OrderBy) Expression() (string, error) {
	key, err := o.getKey()
	if err != nil {
		return "", err
	}
	switch o.Direction {
	case OrderDirectionTypeAscending:
		return fmt.Sprintf("%v ASC", key), nil
	case OrderDirectionTypeDescending:
		return fmt.Sprintf("%v DESC", key), nil
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

func (o *OrderBy) getKey() (string, error) {
	switch o.FieldProperty {
	case FieldPropertyDimension:
		return fmt.Sprintf("%v.%v DESC", o.Table, o.Name), nil
	case FieldPropertyMetric:
		return o.Name, nil
	default:
		return "", fmt.Errorf("not supported field property %v", o.FieldProperty)
	}
}
