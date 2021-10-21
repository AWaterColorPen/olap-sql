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

type DBType string

const(
	DBTypeClickhouse DBType = "DATABASE_TYPE_CLICKHOUSE"
	DBTypeSqlite     DBType = "DATABASE_TYPE_SQLITE"
)

type Column interface {
	GetExpression() string
	GetAlias() string
}

type SingleCol struct {
	Table  string         `json:"table"`
	Name   string         `json:"name"`
	Alias  string         `json:"alias"`
	Type   ColumnType     `json:"type"`
	If     *IfOption      `json:"if"`
	DBType DBType         `json:"dbtype"`
}

type IfCol struct {
	Table string `toml:"table" json:"table"`
	Name  string `toml:"name"  json:"name"`
}

func (i *IfCol) GetExpression() string{
	if i == nil {
		return fmt.Sprintf("NULL")
	}
	return fmt.Sprintf("`%v`.`%v`", i.Table, i.Name)
}

type IfOption struct {
	Filter *Filter `toml:"filter" json:"filter"`
	IfCol1  *IfCol  `toml:"ifcol1"  json:"ifcol1"`
	IfCol2  *IfCol  `toml:"ifcol2"  json:"ifcol2"`
}

func (If *IfOption) GetExpression(dbType DBType) (string, error) {
	filter, err := If.Filter.Expression()
	if err != nil {
		return "", err
	}
	switch dbType {
	case DBTypeClickhouse:
		return fmt.Sprintf("IF(%v,%v,%v)", filter, If.IfCol1.GetExpression(), If.IfCol2.GetExpression()), nil
	case DBTypeSqlite:
		return fmt.Sprintf("IIF(%v, %v, %v)", filter, If.IfCol1.GetExpression(), If.IfCol2.GetExpression()), nil
	default:
		return "", fmt.Errorf("%v unsupport if now", dbType)
	}
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
			If, err := col.If.GetExpression(col.DBType)
			if err != nil {
				return fmt.Sprintf("If error: %v", col.If)
			}
			return fmt.Sprintf("1.0 * COUNT(DISTINCT(%v)) ", If)
		}
		return fmt.Sprintf("1.0 * COUNT(DISTINCT `%v`.`%v` )", col.Table, col.Name)
	case ColumnTypeSum:
		if col.If != nil {
			If, err := col.If.GetExpression(col.DBType)
			if err != nil {
				return fmt.Sprintf("If error: %v", col.If)
			}
			return fmt.Sprintf("1.0 * SUM(%v) ", If)
		}
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
