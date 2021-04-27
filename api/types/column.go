package types

import (
	"fmt"
)

const (
	ColumnTypeCount         ColumnType = "count"
	ColumnTypeDistinctCount ColumnType = "distinct_count"
	ColumnTypeSum           ColumnType = "sum"
	ColumnTypePost          ColumnType = "post"
	ColumnTypeExpression    ColumnType = "expression"
)

type ColumnType string

type Column interface {
	Sql() string
	GetAlias() string
	GetType() ColumnType
}

type CountCol struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (col *CountCol) Sql() string {
	return fmt.Sprintf("COUNT( %s )", col.Name)
}

func (col *CountCol) GetAlias() string {
	return col.Alias
}

func (col *CountCol) GetType() ColumnType {
	return ColumnTypeCount
}

type DistinctCol struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (col *DistinctCol) Sql() string {
	return fmt.Sprintf("1.0 * COUNT(DISTINCT %s )", col.Name)
}

func (col *DistinctCol) GetAlias() string {
	return col.Alias
}

func (col *DistinctCol) GetType() ColumnType {
	return ColumnTypeDistinctCount
}

type SumCol struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (col *SumCol) Sql() string {
	return fmt.Sprintf("1.0 * SUM( %v )", col.Name)
}

func (col *SumCol) GetAlias() string {
	return col.Alias
}

func (col *SumCol) GetType() ColumnType {
	return ColumnTypeSum
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

func (col *ExpressionCol) GetType() ColumnType {
	return ColumnTypeExpression
}
