package types

import (
	"fmt"
	"strings"
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
	DimensionTypeValue      DimensionType = "DIMENSION_VALUE"
	DimensionTypeSingle     DimensionType = "DIMENSION_SINGLE"
	DimensionTypeMulti      DimensionType = "DIMENSION_MULTI"
	DimensionTypeCase       DimensionType = "DIMENSION_CASE"
	DimensionTypeExpression DimensionType = "DIMENSION_EXPRESSION"
)

func (d *Dimension) Expression() (string, error) {
	switch d.Type {
	case DimensionTypeValue:
		return fmt.Sprintf("%v.%v", d.Table, d.Name), nil
	case DimensionTypeSingle:
		return fmt.Sprintf("%v.%v", d.Table, d.FieldName), nil
	case DimensionTypeExpression:
		return d.FieldName, nil
	case DimensionTypeMulti:
		if len(d.Dependency) == 0 {
			return "", fmt.Errorf("dimension dependency len = 0")
		}
		return d.Dependency[0].Expression()
	case DimensionTypeCase:
		if len(d.Dependency) < 2 {
			return "", fmt.Errorf("case dimension dependency len < 2")
		}
		return caseWhenExpression(d.Dependency)
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

func caseWhenExpression(dependencies []*Dimension) (string, error) {
	var statement []string
	for _, v := range dependencies {
		u, err := v.Expression()
		if err != nil {
			return "", err
		}
		// case type can't support integer or float value type
		statement = append(statement, fmt.Sprintf("WHEN %v != '' THEN %v", u, u))
	}
	return fmt.Sprintf("CASE %v END", strings.Join(statement, " ")), nil
}
