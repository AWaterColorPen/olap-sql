package types

import (
	"fmt"
	"strings"
)

type FilterOperatorType string

func (f FilterOperatorType) IsTree() bool {
	switch f {
	case FilterOperatorTypeAnd, FilterOperatorTypeOr:
		return true
	default:
		return false
	}
}

const (
	FilterOperatorTypeUnknown       FilterOperatorType = "FILTER_OPERATOR_UNKNOWN"
	FilterOperatorTypeEquals        FilterOperatorType = "FILTER_OPERATOR_EQUALS"
	FilterOperatorTypeIn            FilterOperatorType = "FILTER_OPERATOR_IN"
	FilterOperatorTypeNotIn         FilterOperatorType = "FILTER_OPERATOR_NOT_IN"
	FilterOperatorTypeLessEquals    FilterOperatorType = "FILTER_OPERATOR_LESS_EQUALS"
	FilterOperatorTypeLess          FilterOperatorType = "FILTER_OPERATOR_LESS"
	FilterOperatorTypeGreaterEquals FilterOperatorType = "FILTER_OPERATOR_GREATER_EQUALS"
	FilterOperatorTypeGreater       FilterOperatorType = "FILTER_OPERATOR_GREATER"
	FilterOperatorTypeLike          FilterOperatorType = "FILTER_OPERATOR_LIKE"
	FilterOperatorTypeExpression    FilterOperatorType = "FILTER_OPERATOR_EXTENSION"
	FilterOperatorTypeAnd           FilterOperatorType = "FILTER_OPERATOR_AND"
	FilterOperatorTypeOr            FilterOperatorType = "FILTER_OPERATOR_OR"
)

type ValueType string

const (
	ValueTypeUnknown ValueType = "VALUE_UNKNOWN"
	ValueTypeString  ValueType = "VALUE_STRING"
	ValueTypeInteger ValueType = "VALUE_INTEGER"
	ValueTypeFloat   ValueType = "VALUE_FLOAT"
)

type Filter struct {
	OperatorType FilterOperatorType `toml:"operator_type" json:"operator_type"`
	ValueType    ValueType          `toml:"value_type"    json:"value_type"`
	Name         string             `toml:"name"          json:"name"`
	Value        []interface{}      `toml:"value"         json:"value"`
	Children     []*Filter          `toml:"children"      json:"children"`
}

func (f *Filter) Expression() (string, error) {
	value, err := f.valueToStringSlice()
	if err != nil {
		return "", err
	}
	switch f.OperatorType {
	case FilterOperatorTypeEquals:
		v := value[0]
		return fmt.Sprintf("%v = %v", f.Name, v), nil
	case FilterOperatorTypeIn:
		return fmt.Sprintf("%v IN (%v)", f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorTypeNotIn:
		return fmt.Sprintf("%v NOT IN (%v)", f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorTypeLessEquals:
		v := value[0]
		return fmt.Sprintf("%v <= %v", f.Name, v), nil
	case FilterOperatorTypeLess:
		v := value[0]
		return fmt.Sprintf("%v < %v", f.Name, v), nil
	case FilterOperatorTypeGreaterEquals:
		v := value[0]
		return fmt.Sprintf("%v >= %v", f.Name, v), nil
	case FilterOperatorTypeGreater:
		v := value[0]
		return fmt.Sprintf("%v > %v", f.Name, v), nil
	case FilterOperatorTypeLike:
		v := value[0]
		return fmt.Sprintf("%v LIKE %v", f.Name, v), nil
	case FilterOperatorTypeExpression:
		v := value[0]
		return v, nil
	case FilterOperatorTypeAnd:
		return f.treeStatement(" AND ")
	case FilterOperatorTypeOr:
		return f.treeStatement(" OR ")
	default:
		return "", fmt.Errorf("not supported filter operator type %v", f.OperatorType)
	}
}

func (f *Filter) Alias() (string, error) {
	return "", fmt.Errorf("filter is unsupported alias method")
}

func (f *Filter) Statement() (string, error) {
	return f.Expression()
}

func (f *Filter) valueToStringSlice() ([]string, error) {
	var out []string
	for _, v := range f.Value {
		switch f.ValueType {
		case ValueTypeString:
			out = append(out, fmt.Sprintf("'%v'", v))
		case ValueTypeInteger, ValueTypeFloat:
			out = append(out, fmt.Sprintf("%v", v))
		default:
			return nil, fmt.Errorf("not supported value type %v", f.ValueType)
		}
	}
	return out, nil
}

func (f *Filter) treeStatement(sep string) (string, error) {
	var filter []string
	for _, v := range f.Children {
		statement, err := v.Statement()
		if err != nil {
			return "", err
		}
		filter = append(filter, statement)
	}
	return fmt.Sprintf("( %v )", strings.Join(filter, sep)), nil
}
