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
	Table       string        `json:"table"`
	Name        string        `json:"name"`
	FieldName   string        `json:"field_name"`
	Type        DimensionType `json:"type"`
	Composition []string      `json:"composition"`
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
		if len(d.Composition) == 0 {
			return "", fmt.Errorf("dimension composition error")
		}
		return fmt.Sprintf("%v", d.Composition[0]), nil
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

func (d *Dimension) expression() bool {
	return !reg.MatchString(d.FieldName)
}
