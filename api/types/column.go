package types

import (
	"fmt"
	"strings"
)

const (
	ColumnTypeValue         ColumnType = "value"
	ColumnTypeCount         ColumnType = "count"
	ColumnTypeDistinctCount ColumnType = "distinct_count"
	ColumnTypeSum           ColumnType = "sum"
	ColumnTypeAdd           ColumnType = "add"
	ColumnTypeSubtract      ColumnType = "subtract"
	ColumnTypeMultiply      ColumnType = "multiply"
	ColumnTypeDivide        ColumnType = "divide"
	ColumnTypePost          ColumnType = "post"
	ColumnTypeExpression    ColumnType = "expression"
)

type ColumnType string

type Column interface {
	GetExpression() string
	GetAlias() string
}

type SingleCol struct {
	Table string     `json:"table"`
	Name  string     `json:"name"`
	Alias string     `json:"alias"`
	Type  ColumnType `json:"type"`
	If	  *IfOption   `json:"if"`
}

type IfOption struct {
	Filter  *Filter `json:"filter"`
	Metric1 *Metric `json:"metric1"`
	Metric2 *Metric `json:"metric2"`
}

func (If *IfOption) GetExpression() string {
	filter, _ := If.Filter.Expression()
	metric1, _ := If.Metric1.Expression()
	metric2, _ := If.Metric2.Expression()

	return fmt.Sprintf("if(`%v`, `%v`, `%v`)",filter, metric1, metric2)
}

func (col *SingleCol) GetExpression() string {
	switch col.Type {
	case ColumnTypeValue:
		return fmt.Sprintf("`%v`.`%v`", col.Table, col.Name)
	case ColumnTypeCount:
		if col.Name == "*" {
			return fmt.Sprintf("COUNT(*)")
		}
		return fmt.Sprintf("COUNT( `%v`.`%v` )", col.Table, col.Name)
	case ColumnTypeDistinctCount:
		if col.If != nil {
			return fmt.Sprintf("1.0 * COUNT(DISTINCT(`%v`) ", col.If.GetExpression())
		}
		return fmt.Sprintf("1.0 * COUNT(DISTINCT `%v`.`%v` )", col.Table, col.Name)
	case ColumnTypeSum:
		return fmt.Sprintf("1.0 * SUM( `%v`.`%v` )", col.Table, col.Name)
	default:
		return fmt.Sprintf("unsupported type: %v", col.Type)
	}
}

func (col *SingleCol) GetAlias() string {
	return col.Alias
}

type ArithmeticOperatorType string

const (
	ArithmeticOperatorTypeAdd      ArithmeticOperatorType = "+"
	ArithmeticOperatorTypeSubtract ArithmeticOperatorType = "-"
	ArithmeticOperatorTypeMultiply ArithmeticOperatorType = "*"
	ArithmeticOperatorTypeDivide   ArithmeticOperatorType = "/"
)

type ArithmeticCol struct {
	Column []Column   `json:"column"`
	Alias  string     `json:"alias"`
	Type   ColumnType `json:"type"`
}

func (col *ArithmeticCol) GetExpression() string {
	var children []string
	for _, v := range col.Column {
		children = append(children, fmt.Sprintf("( %v )", v.GetExpression()))
	}
	operator := ArithmeticOperatorType("")
	switch col.Type {
	case ColumnTypeAdd:
		operator = ArithmeticOperatorTypeAdd
	case ColumnTypeSubtract:
		operator = ArithmeticOperatorTypeSubtract
	case ColumnTypeMultiply:
		operator = ArithmeticOperatorTypeMultiply
	case ColumnTypeDivide:
		operator = ArithmeticOperatorTypeDivide
	}
	if operator == ArithmeticOperatorTypeDivide {
		for i := 1; i < len(children); i++ {
			children[i] = fmt.Sprintf("NULLIF(%v, 0)", children[i])
		}
	}
	return fmt.Sprintf("( %v )", strings.Join(children, fmt.Sprintf(" %v  ", operator)))
}

func (col *ArithmeticCol) GetAlias() string {
	return col.Alias
}

type ExpressionCol struct {
	Expression string `json:"expression"`
	Alias      string `json:"alias"`
}

func (col *ExpressionCol) GetExpression() string {
	return col.Expression
}

func (col *ExpressionCol) GetAlias() string {
	return col.Alias
}
