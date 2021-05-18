package types

import (
	"fmt"
	"strings"

	"github.com/awatercolorpen/olap-sql/api/proto"
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
	FilterOperatorTypeUnknown       FilterOperatorType = "FILTER_OPERATOR_UNKNOWN"
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

func (v ValueType) ToEnum() proto.VALUE_TYPE {
	if u, ok := proto.VALUE_TYPE_value[string(v)]; ok {
		return proto.VALUE_TYPE(u)
	}
	return proto.VALUE_TYPE_VALUE_UNKNOWN
}

func EnumToValueType(v proto.VALUE_TYPE) ValueType {
	return ValueType(v.String())
}

const (
	ValueTypeUnknown ValueType = "VALUE_UNKNOWN"
	ValueTypeString  ValueType = "VALUE_STRING"
	ValueTypeInteger ValueType = "VALUE_INTEGER"
	ValueTypeFloat   ValueType = "VALUE_FLOAT"
)

type Filter struct {
	OperatorType FilterOperatorType `json:"operator_type"`
	ValueType    ValueType          `json:"value_type"`
	Table        string             `json:"table"`
	Name         string             `json:"name"`
	Value        []interface{}      `json:"value"`
	Children     []*Filter          `json:"children"`
}

func (f *Filter) Statement() (string, error) {
	value, err := f.valueToStringSlice()
	if err != nil {
		return "", err
	}
	switch f.OperatorType {
	case FilterOperatorTypeIn:
		return fmt.Sprintf("`%v`.`%v` IN (%v)", f.Table, f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorTypeNotIn:
		return fmt.Sprintf("`%v`.`%v` NOT IN (%v)", f.Table, f.Name, strings.Join(value, ", ")), nil
	case FilterOperatorTypeLessEquals:
		v := value[0]
		return fmt.Sprintf("`%v`.`%v` <= %v", f.Table, f.Name, v), nil
	case FilterOperatorTypeLess:
		v := value[0]
		return fmt.Sprintf("`%v`.`%v` < %v", f.Table, f.Name, v), nil
	case FilterOperatorTypeGreaterEquals:
		v := value[0]
		return fmt.Sprintf("`%v`.`%v` >= %v", f.Table, f.Name, v), nil
	case FilterOperatorTypeGreater:
		v := value[0]
		return fmt.Sprintf("`%v`.`%v` > %v", f.Table, f.Name, v), nil
	case FilterOperatorTypeLike:
		v := value[0]
		return fmt.Sprintf("`%v`.`%v` LIKE %v", f.Table, f.Name, v), nil
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

func (f *Filter) ToProto() *proto.Filter {
	return &proto.Filter{

	}
}

func ProtoToFilter(m *proto.Filter) *Filter {
	return &Filter{}
}
