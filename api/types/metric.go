package types

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

type MetricType string

func (m MetricType) ToEnum() proto.METRIC_TYPE {
	if v, ok := proto.METRIC_TYPE_value[string(m)]; ok {
		return proto.METRIC_TYPE(v)
	}
	return proto.METRIC_TYPE_METRIC_TYPE_UNKNOWN
}

func EnumToMetricType(m proto.METRIC_TYPE) MetricType {
	return MetricType(m.String())
}

const (
	MetricTypeUnknown       MetricType = "METRIC_TYPE_UNKNOWN"
	MetricTypeValue         MetricType = "METRIC_TYPE_VALUE"
	MetricTypeCount         MetricType = "METRIC_TYPE_COUNT"
	MetricTypeDistinctCount MetricType = "METRIC_TYPE_DISTINCT_COUNT"
	MetricTypeSum           MetricType = "METRIC_TYPE_SUM"
	MetricTypeAdd           MetricType = "METRIC_TYPE_ADD"
	MetricTypeSubtract      MetricType = "METRIC_TYPE_SUBTRACT"
	MetricTypeMultiply      MetricType = "METRIC_TYPE_MULTIPLY"
	MetricTypeDivide        MetricType = "METRIC_TYPE_DIVIDE"
	MetricTypePost          MetricType = "METRIC_TYPE_POST"
	MetricTypeExpression    MetricType = "METRIC_TYPE_EXTENSION"
)

type Metric struct {
	Type           MetricType  `json:"type"`
	Table          string      `json:"table"`
	Name           string      `json:"name"`
	FieldName      string      `json:"field_name"`
	ExtensionValue interface{} `json:"extension_value"`
}

func (m *Metric) Statement() (string, error) {
	col, err := m.column()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v AS %v", col.Sql(), col.GetAlias()), nil
}

func (m *Metric) column() (Column, error) {
	switch m.Type {
	case MetricTypeValue:
		return &SingleCol{Table: m.Table, Name: m.FieldName, Alias: m.Name, Type: ColumnTypeValue}, nil
	case MetricTypeCount:
		return &SingleCol{Table: m.Table, Name: m.FieldName, Alias: m.Name, Type: ColumnTypeCount}, nil
	case MetricTypeDistinctCount:
		return &SingleCol{Table: m.Table, Name: m.FieldName, Alias: m.Name, Type: ColumnTypeDistinctCount}, nil
	case MetricTypeSum:
		return &SingleCol{Table: m.Table, Name: m.FieldName, Alias: m.Name, Type: ColumnTypeMultiply}, nil
	case MetricTypeAdd:
		// todo column
		return &ArithmeticCol{Column: []Column{}, Alias: m.Name, Type: ColumnTypeAdd}, nil
	case MetricTypeSubtract:
		return &ArithmeticCol{Column: []Column{}, Alias: m.Name, Type: ColumnTypeSubtract}, nil
	case MetricTypeMultiply:
		return &ArithmeticCol{Column: []Column{}, Alias: m.Name, Type: ColumnTypeMultiply}, nil
	case MetricTypeDivide:
		return &ArithmeticCol{Column: []Column{}, Alias: m.Name, Type: ColumnTypeDivide}, nil
	case MetricTypeExpression:
		expression, ok := m.ExtensionValue.(string)
		if !ok {
			return nil, fmt.Errorf("invalid metric type expression %v", m.ExtensionValue)
		}
		return &ExpressionCol{Expression: expression, Alias: m.Name}, nil
	default:
		return nil, fmt.Errorf("not supported metric type %v", m.Type)
	}
}

func (m *Metric) ToProto() *proto.Metric {
	return &proto.Metric{
		Type:           m.Type.ToEnum(),
		Table:          m.Table,
		Name:           m.Name,
		FieldName:      m.FieldName,
		ExtensionValue: fmt.Sprint(m.ExtensionValue),
	}
}

func ProtoToMetric(m *proto.Metric) *Metric {
	return &Metric{}
}
