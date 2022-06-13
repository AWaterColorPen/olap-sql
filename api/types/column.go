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
	ColumnTypeAs            ColumnType = "as"
	ColumnTypePost          ColumnType = "post"
	ColumnTypeExpression    ColumnType = "expression"
)

type ColumnType string

type Column interface {
	GetExpression() string
	GetAlias() string
}

type SingleCol struct {
	DBType DBType     `json:"db_type"`
	Table  string     `json:"table"`
	Name   string     `json:"name"`
	Alias  string     `json:"alias"`
	Type   ColumnType `json:"type"`
	Filter *Filter    `json:"filter"`
}

func (col *SingleCol) GetIfExpression() (string, error) {
	filter, err := col.Filter.Expression()
	if err != nil {
		return "", err
	}
	switch col.DBType {
	case DBTypeClickHouse:
		return fmt.Sprintf("IF( %v , `%v`.`%v` , NULL )", filter, col.Table, col.Name), nil
	case DBTypeSQLite:
		return fmt.Sprintf("IIF( %v , `%v`.`%v` , NULL )", filter, col.Table, col.Name), nil
	default:
		return "", fmt.Errorf("%v unsupport if now", col.DBType)
	}
}

func (col *SingleCol) GetExpression() string {
	switch col.Type {
	case ColumnTypeValue:
		return col.getSimpleName()
	case ColumnTypeCount:
		if col.Name == "*" {
			return fmt.Sprintf("COUNT(*)")
		}
		return fmt.Sprintf("COUNT( `%v`.`%v` )", col.Table, col.Name)
	case ColumnTypeDistinctCount:
		name := col.getSimpleName()
		var err error
		if col.Filter != nil {
			name, err = col.GetIfExpression()
			if err != nil {
				return fmt.Sprintf("%v", err)
			}
		}
		return fmt.Sprintf("1.0 * COUNT(DISTINCT %v )", name)
	case ColumnTypeSum:
		name := col.getSimpleName()
		var err error
		if col.Filter != nil {
			name, err = col.GetIfExpression()
			if err != nil {
				return fmt.Sprintf("%v", err)
			}
		}
		return fmt.Sprintf("1.0 * SUM(%v)", name)
	default:
		return fmt.Sprintf("unsupported type: %v", col.Type)
	}
}

func (col *SingleCol) GetAlias() string {
	return col.Alias
}

func (col *SingleCol) getSimpleName() string {
	if col.Name == "" {
		return fmt.Sprintf("`%v`.`%v`", col.Table, col.Alias)
	}
	return fmt.Sprintf("`%v`.`%v`", col.Table, col.Name)
}

type ArithmeticOperatorType string

const (
	ArithmeticOperatorTypeAdd      ArithmeticOperatorType = "+"
	ArithmeticOperatorTypeSubtract ArithmeticOperatorType = "-"
	ArithmeticOperatorTypeMultiply ArithmeticOperatorType = "*"
	ArithmeticOperatorTypeDivide   ArithmeticOperatorType = "/"
	ArithmeticOperatorTypeAs       ArithmeticOperatorType = "as"
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
	case ColumnTypeAs:
		operator = ArithmeticOperatorTypeAs
	}
	if operator == ArithmeticOperatorTypeAs {
		return fmt.Sprintf("( %v )", strings.Join(children, ""))
	}
	if operator == ArithmeticOperatorTypeDivide {
		for i := 1; i < len(children); i++ {
			children[i] = fmt.Sprintf("NULLIF(%v, 0)", children[i])
		}
	}
	if operator == ArithmeticOperatorTypeAdd {
		for i := 1; i < len(children); i++ {
			children[i] = fmt.Sprintf("IFNULL(%v, 0)", children[i])
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
