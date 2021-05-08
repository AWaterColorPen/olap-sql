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
	Sql() string
	GetAlias() string
}

type SingleCol struct {
	Table string     `json:"table"`
	Name  string     `json:"name"`
	Alias string     `json:"alias"`
	Type  ColumnType `json:"type"`
}

func (col *SingleCol) Sql() string {
	switch col.Type {
	case ColumnTypeValue:
		return fmt.Sprintf("`%v`.`%v`", col.Table, col.Name)
	case ColumnTypeCount:
		if col.Name == "*" {
			return fmt.Sprintf("COUNT(*)")
		}
		return fmt.Sprintf("COUNT( `%v`.`%v` )", col.Table, col.Name)
	case ColumnTypeDistinctCount:
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

func (col *ArithmeticCol) Sql() string {
	var son []string
	for _, v := range col.Column {
		son = append(son, fmt.Sprintf("( %v )", v.Sql()))
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
	return fmt.Sprintf("( %v )", strings.Join(son, fmt.Sprintf(" %v  ", operator)))
}

func (col *ArithmeticCol) GetAlias() string {
	return col.Alias
}

type ExpressionCol struct {
	Expression string `json:"expression"`
	Alias      string `json:"alias"`
}

func (col *ExpressionCol) Sql() string {
	return col.Expression
}

func (col *ExpressionCol) GetAlias() string {
	return col.Alias
}
