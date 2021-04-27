package types

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

var (
	ErrNotSupportedMetricType      = fmt.Errorf("not supported MetricType")
	ErrInvalidMetricTypeExpression = fmt.Errorf("invalid metric type expression")
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
	MetricTypeCount         MetricType = "METRIC_TYPE_COUNT"
	MetricTypeDistinctCount MetricType = "METRIC_TYPE_DISTINCT_COUNT"
	MetricTypeSum           MetricType = "METRIC_TYPE_SUM"
	MetricTypePost          MetricType = "METRIC_TYPE_POST"
	MetricTypeExpression    MetricType = "METRIC_TYPE_EXTENSION"
)

type Metric struct {
	Type           MetricType  `json:"type"`
	FieldName      string      `json:"field_name"`
	Name           string      `json:"name"`
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
	case MetricTypeCount:
		return &CountCol{Name: m.FieldName, Alias: m.Name}, nil
	case MetricTypeDistinctCount:
		return &DistinctCol{Name: m.FieldName, Alias: m.Name}, nil
	case MetricTypeSum:
		return &SumCol{Name: m.FieldName, Alias: m.Name}, nil
	case MetricTypeExpression:
		expression, ok := m.ExtensionValue.(string)
		if !ok {
			return nil, ErrInvalidMetricTypeExpression
		}
		return &ExpressionCol{Expression: expression, Alias: m.Name}, nil
	default:
		return nil, ErrNotSupportedMetricType
	}
}

func (m *Metric) ToProto() *proto.Metric {
	return &proto.Metric{
		Type:           m.Type.ToEnum(),
		FieldName:      m.FieldName,
		Name:           m.Name,
		ExtensionValue: fmt.Sprint(m.ExtensionValue),
	}
}

func ProtoToMetric(m *proto.Metric) *Metric {
	return &Metric{}
}
