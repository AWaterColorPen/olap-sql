package types

import (
	"fmt"
	"strings"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

var (
	ErrNotSupportedFilterOperatorType = fmt.Errorf("not supported FilterOperatorType")
	ErrNotSupportedFilterValueType    = fmt.Errorf("not supported FilterValueType")
)

type FilterOperatorType string

func (f FilterOperatorType) ToEnum() proto.FILTER_OPERATOR_TYPE {
	if v, ok := proto.FILTER_OPERATOR_TYPE_value[string(f)]; ok {
		return proto.FILTER_OPERATOR_TYPE(v)
	}
	return proto.FILTER_OPERATOR_TYPE_FILTER_OPERATOR_UNKNOWN
}

func EnumToFilterOperatorType(f proto.FILTER_OPERATOR_TYPE) FilterOperatorType {
	return FilterOperatorType(f.String())
}

const (
	FilterOperatorUnknown       FilterOperatorType = "FILTER_OPERATOR_UNKNOWN"
	FilterOperatorIn            FilterOperatorType = "FILTER_OPERATOR_IN"
	FilterOperatorNotIn         FilterOperatorType = "FILTER_OPERATOR_NOT_IN"
	FilterOperatorLessEquals    FilterOperatorType = "FILTER_OPERATOR_LESS_EQUALS"
	FilterOperatorLess          FilterOperatorType = "FILTER_OPERATOR_LESS"
	FilterOperatorGreaterEquals FilterOperatorType = "FILTER_OPERATOR_GREATER_EQUALS"
	FilterOperatorGreater       FilterOperatorType = "FILTER_OPERATOR_GREATER"
	FilterOperatorLike          FilterOperatorType = "FILTER_OPERATOR_LIKE"
	FilterOperatorExpression    FilterOperatorType = "FILTER_OPERATOR_EXTENSION"
)

type FilterValueType string

func (f FilterValueType) ToEnum() proto.FILTER_VALUE_TYPE {
	if v, ok := proto.FILTER_VALUE_TYPE_value[string(f)]; ok {
		return proto.FILTER_VALUE_TYPE(v)
	}
	return proto.FILTER_VALUE_TYPE_FILTER_VALUE_UNKNOWN
}

func EnumToFilterValueType(f proto.FILTER_VALUE_TYPE) FilterValueType {
	return FilterValueType(f.String())
}

const (
	FilterValueUnknown FilterValueType = "FILTER_VALUE_UNKNOWN"
	FilterValueString  FilterValueType = "FILTER_VALUE_STRING"
	FilterValueInteger FilterValueType = "FILTER_VALUE_INTEGER"
	FilterValueFloat   FilterValueType = "FILTER_VALUE_FLOAT"
)

type Filter struct {
	OperatorType FilterOperatorType `json:"operator_type"`
	ValueType    FilterValueType    `json:"value_type"`
	Name         string             `json:"name"`
	Value        []interface{}      `json:"value"`
}

type Filters []*Filter

func (f *Filter) Statement() (string, error) {
	value, err := f.valueToStringSlice()
	if err != nil {
		return "", err
	}
	switch f.OperatorType {
	case FilterOperatorIn:
		return fmt.Sprintf("%v IN (%v)", f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorNotIn:
		return fmt.Sprintf("%v NOT IN (%v)", f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorLessEquals:
		v := value[0]
		return fmt.Sprintf("%v <= %v", f.Name, v), nil
	case FilterOperatorLess:
		v := value[0]
		return fmt.Sprintf("%v < %v", f.Name, v), nil
	case FilterOperatorGreaterEquals:
		v := value[0]
		return fmt.Sprintf("%v >= %v", f.Name, v), nil
	case FilterOperatorGreater:
		v := value[0]
		return fmt.Sprintf("%v > %v", f.Name, v), nil
	case FilterOperatorLike:
		v := value[0]
		return fmt.Sprintf("%v LIKE %v", f.Name, v), nil
	case FilterOperatorExpression:
		v := value[0]
		return v, nil
	default:
		return "", ErrNotSupportedFilterOperatorType
	}
}

func (f *Filter) valueToStringSlice() ([]string, error) {
	var out []string
	for _, v := range f.Value {
		switch f.ValueType {
		case FilterValueString:
			out = append(out, fmt.Sprintf("'%v'", v))
		case FilterValueInteger, FilterValueFloat:
			out = append(out, fmt.Sprintf("%v", v))
		default:
			return nil, ErrNotSupportedFilterValueType
		}
	}
	return out, nil
}
