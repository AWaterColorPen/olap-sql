package types

import (
	"fmt"
	"regexp"
)

var (
	reg = regexp.MustCompile(`^[0-9A-Za-z_]+$`)
)

type DimensionType string

type Dimension struct {
	Table      string        `json:"table"`
	Name       string        `json:"name"`
	FieldName  string        `json:"field_name"`
	Type       DimensionType `json:"type"`
	ValueType  ValueType     `json:"value_type"`
	Dependency []*Dimension  `json:"dependency"`
}

const (
	DimensionTypeSingle     DimensionType = "DIMENSION_SINGLE"
	DimensionTypeMulti      DimensionType = "DIMENSION_MULTI"
	DimensionTypeExpression DimensionType = "DIMENSION_EXPRESSION"
)

func (d *Dimension) Expression() (string, error) {
	switch d.Type {
	case DimensionTypeSingle:
		return fmt.Sprintf("%v.%v", d.Table, d.FieldName), nil
	case DimensionTypeExpression:
		return d.FieldName, nil
	case DimensionTypeMulti:
		if len(d.Dependency) == 0 {
			return "", fmt.Errorf("dimension dependency error")
		}
		return d.Dependency[0].Expression()
	default:
		return "", fmt.Errorf("unsupported dimension type")
	}
}

func (d *Dimension) Alias() (string, error) {
	return d.Name, nil
}

func (d *Dimension) Statement() (string, error) {
	expression, err := d.Expression()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v AS %v", expression, d.Name), nil
}
