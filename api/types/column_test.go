package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleCol_GetExpression_Value(t *testing.T) {
	col := &SingleCol{
		DBType: DBTypeSQLite,
		Table:  "orders",
		Name:   "status",
		Alias:  "order_status",
		Type:   ColumnTypeValue,
	}
	assert.Equal(t, "`orders`.`status`", col.GetExpression())
	assert.Equal(t, "order_status", col.GetAlias())
}

func TestSingleCol_GetExpression_Count(t *testing.T) {
	t.Run("count star", func(t *testing.T) {
		col := &SingleCol{
			Table: "orders",
			Name:  "*",
			Alias: "cnt",
			Type:  ColumnTypeCount,
		}
		assert.Equal(t, "COUNT(*)", col.GetExpression())
	})

	t.Run("count field", func(t *testing.T) {
		col := &SingleCol{
			Table: "orders",
			Name:  "id",
			Alias: "cnt",
			Type:  ColumnTypeCount,
		}
		assert.Equal(t, "COUNT( `orders`.`id` )", col.GetExpression())
	})
}

func TestSingleCol_GetExpression_DistinctCount(t *testing.T) {
	t.Run("no filter", func(t *testing.T) {
		col := &SingleCol{
			Table: "orders",
			Name:  "user_id",
			Alias: "uv",
			Type:  ColumnTypeDistinctCount,
		}
		assert.Equal(t, "1.0 * COUNT(DISTINCT `orders`.`user_id` )", col.GetExpression())
	})

	t.Run("with clickhouse filter", func(t *testing.T) {
		filter := &Filter{
			OperatorType: FilterOperatorTypeEquals,
			ValueType:    ValueTypeString,
			Name:         "type",
			Value:        []any{"vip"},
		}
		col := &SingleCol{
			DBType: DBTypeClickHouse,
			Table:  "orders",
			Name:   "user_id",
			Alias:  "vip_uv",
			Type:   ColumnTypeDistinctCount,
			Filter: filter,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "COUNT(DISTINCT")
		assert.Contains(t, expr, "IF(")
	})
}

func TestSingleCol_GetExpression_Sum(t *testing.T) {
	t.Run("no filter", func(t *testing.T) {
		col := &SingleCol{
			Table: "orders",
			Name:  "amount",
			Alias: "total",
			Type:  ColumnTypeSum,
		}
		assert.Equal(t, "1.0 * SUM(`orders`.`amount`)", col.GetExpression())
	})

	t.Run("with sqlite filter", func(t *testing.T) {
		filter := &Filter{
			OperatorType: FilterOperatorTypeEquals,
			ValueType:    ValueTypeString,
			Name:         "status",
			Value:        []any{"paid"},
		}
		col := &SingleCol{
			DBType: DBTypeSQLite,
			Table:  "orders",
			Name:   "amount",
			Alias:  "paid_amount",
			Type:   ColumnTypeSum,
			Filter: filter,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "SUM(")
		assert.Contains(t, expr, "IIF(")
	})
}

func TestSingleCol_GetExpression_Unsupported(t *testing.T) {
	col := &SingleCol{
		Table: "t",
		Name:  "x",
		Type:  "unknown_type",
	}
	expr := col.GetExpression()
	assert.Contains(t, expr, "unsupported type")
}

func TestSingleCol_GetExpression_EmptyName(t *testing.T) {
	// When Name is empty, getSimpleName uses Alias
	col := &SingleCol{
		Table: "orders",
		Name:  "",
		Alias: "my_alias",
		Type:  ColumnTypeValue,
	}
	assert.Equal(t, "`orders`.`my_alias`", col.GetExpression())
}

func TestSingleCol_GetIfExpression(t *testing.T) {
	filter := &Filter{
		OperatorType: FilterOperatorTypeEquals,
		ValueType:    ValueTypeString,
		Name:         "status",
		Value:        []any{"active"},
	}

	t.Run("clickhouse", func(t *testing.T) {
		col := &SingleCol{
			DBType: DBTypeClickHouse,
			Table:  "users",
			Name:   "id",
			Filter: filter,
		}
		expr, err := col.GetIfExpression()
		require.NoError(t, err)
		assert.Contains(t, expr, "IF(")
		assert.Contains(t, expr, "`users`.`id`")
	})

	t.Run("sqlite", func(t *testing.T) {
		col := &SingleCol{
			DBType: DBTypeSQLite,
			Table:  "users",
			Name:   "id",
			Filter: filter,
		}
		expr, err := col.GetIfExpression()
		require.NoError(t, err)
		assert.Contains(t, expr, "IIF(")
	})

	t.Run("unsupported dbtype", func(t *testing.T) {
		col := &SingleCol{
			DBType: "unknown_db",
			Table:  "users",
			Name:   "id",
			Filter: filter,
		}
		_, err := col.GetIfExpression()
		assert.Error(t, err)
	})
}

func TestArithmeticCol_GetExpression(t *testing.T) {
	left := &SingleCol{Table: "t", Name: "a", Alias: "a", Type: ColumnTypeValue}
	right := &SingleCol{Table: "t", Name: "b", Alias: "b", Type: ColumnTypeValue}

	t.Run("add", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left, right},
			Alias:  "sum_ab",
			Type:   ColumnTypeAdd,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "+")
		assert.Contains(t, expr, "IFNULL")
	})

	t.Run("subtract", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left, right},
			Alias:  "diff_ab",
			Type:   ColumnTypeSubtract,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "-")
	})

	t.Run("multiply", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left, right},
			Alias:  "mul_ab",
			Type:   ColumnTypeMultiply,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "*")
	})

	t.Run("divide", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left, right},
			Alias:  "div_ab",
			Type:   ColumnTypeDivide,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "/")
		assert.Contains(t, expr, "NULLIF")
	})

	t.Run("as", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left},
			Alias:  "as_a",
			Type:   ColumnTypeAs,
		}
		expr := col.GetExpression()
		assert.Contains(t, expr, "`t`.`a`")
	})

	t.Run("alias", func(t *testing.T) {
		col := &ArithmeticCol{
			Column: []Column{left, right},
			Alias:  "result",
			Type:   ColumnTypeAdd,
		}
		assert.Equal(t, "result", col.GetAlias())
	})
}

func TestExpressionCol(t *testing.T) {
	col := &ExpressionCol{
		Expression: "CUSTOM_FUNC(x, y)",
		Alias:      "custom",
	}
	assert.Equal(t, "CUSTOM_FUNC(x, y)", col.GetExpression())
	assert.Equal(t, "custom", col.GetAlias())
}
