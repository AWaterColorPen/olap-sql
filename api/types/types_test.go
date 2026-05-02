package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- Dimension ----

func TestDimension_Expression(t *testing.T) {
	tests := []struct {
		name      string
		dim       *Dimension
		expected  string
		expectErr bool
	}{
		{
			name:     "value type",
			dim:      &Dimension{Table: "orders", Name: "city", Type: DimensionTypeValue},
			expected: "orders.city",
		},
		{
			name:     "single type",
			dim:      &Dimension{Table: "users", Name: "country", FieldName: "country_code", Type: DimensionTypeSingle},
			expected: "users.country_code",
		},
		{
			name:     "expression type",
			dim:      &Dimension{Name: "custom", FieldName: "DATE(created_at)", Type: DimensionTypeExpression},
			expected: "DATE(created_at)",
		},
		{
			name: "multi type",
			dim: &Dimension{
				Name: "multi",
				Type: DimensionTypeMulti,
				Dependency: []*Dimension{
					{Table: "t1", Name: "col", Type: DimensionTypeValue},
				},
			},
			expected: "t1.col",
		},
		{
			name:      "multi type empty dependency",
			dim:       &Dimension{Name: "multi", Type: DimensionTypeMulti, Dependency: nil},
			expectErr: true,
		},
		{
			name: "case type",
			dim: &Dimension{
				Name: "case_dim",
				Type: DimensionTypeCase,
				Dependency: []*Dimension{
					{Table: "t1", Name: "a", Type: DimensionTypeValue},
					{Table: "t2", Name: "b", Type: DimensionTypeValue},
				},
			},
			expected: "CASE WHEN t1.a != '' THEN t1.a WHEN t2.b != '' THEN t2.b END",
		},
		{
			name:      "case type too few dependencies",
			dim:       &Dimension{Name: "c", Type: DimensionTypeCase, Dependency: []*Dimension{{Table: "t", Name: "x", Type: DimensionTypeValue}}},
			expectErr: true,
		},
		{
			name:      "unsupported type",
			dim:       &Dimension{Name: "x", Type: "UNSUPPORTED"},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dim.Expression()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestDimension_Alias(t *testing.T) {
	d := &Dimension{Name: "city"}
	alias, err := d.Alias()
	require.NoError(t, err)
	assert.Equal(t, "city", alias)
}

func TestDimension_Statement(t *testing.T) {
	d := &Dimension{Table: "orders", Name: "city", Type: DimensionTypeValue}
	stmt, err := d.Statement()
	require.NoError(t, err)
	assert.Equal(t, "orders.city AS city", stmt)
}

// ---- Metric ----

func TestMetric_Expression(t *testing.T) {
	tests := []struct {
		name      string
		metric    *Metric
		contains  string
		expectErr bool
	}{
		{
			name:     "value",
			metric:   &Metric{Table: "orders", Name: "status", FieldName: "status", Type: MetricTypeValue},
			contains: "`orders`.`status`",
		},
		{
			name:     "count",
			metric:   &Metric{Table: "orders", Name: "cnt", FieldName: "id", Type: MetricTypeCount},
			contains: "COUNT( `orders`.`id` )",
		},
		{
			name:     "count star",
			metric:   &Metric{Table: "orders", Name: "cnt", FieldName: "*", Type: MetricTypeCount},
			contains: "COUNT(*)",
		},
		{
			name:     "distinct_count",
			metric:   &Metric{Table: "orders", Name: "uv", FieldName: "user_id", Type: MetricTypeDistinctCount},
			contains: "COUNT(DISTINCT",
		},
		{
			name:     "sum",
			metric:   &Metric{Table: "orders", Name: "total", FieldName: "amount", Type: MetricTypeSum},
			contains: "SUM(",
		},
		{
			name:     "expression",
			metric:   &Metric{Name: "custom", FieldName: "CUSTOM_FUNC(x)", Type: MetricTypeExpression},
			contains: "CUSTOM_FUNC(x)",
		},
		{
			name: "add",
			metric: &Metric{
				Name: "ab",
				Type: MetricTypeAdd,
				Children: []*Metric{
					{Table: "t", Name: "a", FieldName: "a", Type: MetricTypeValue},
					{Table: "t", Name: "b", FieldName: "b", Type: MetricTypeValue},
				},
			},
			contains: "+",
		},
		{
			name: "subtract",
			metric: &Metric{
				Name: "diff",
				Type: MetricTypeSubtract,
				Children: []*Metric{
					{Table: "t", Name: "a", FieldName: "a", Type: MetricTypeValue},
					{Table: "t", Name: "b", FieldName: "b", Type: MetricTypeValue},
				},
			},
			contains: "-",
		},
		{
			name: "multiply",
			metric: &Metric{
				Name: "mul",
				Type: MetricTypeMultiply,
				Children: []*Metric{
					{Table: "t", Name: "a", FieldName: "a", Type: MetricTypeValue},
					{Table: "t", Name: "b", FieldName: "b", Type: MetricTypeValue},
				},
			},
			contains: "*",
		},
		{
			name: "divide",
			metric: &Metric{
				Name: "div",
				Type: MetricTypeDivide,
				Children: []*Metric{
					{Table: "t", Name: "a", FieldName: "a", Type: MetricTypeValue},
					{Table: "t", Name: "b", FieldName: "b", Type: MetricTypeValue},
				},
			},
			contains: "NULLIF",
		},
		{
			name: "as",
			metric: &Metric{
				Name: "as_m",
				Type: MetricTypeAs,
				Children: []*Metric{
					{Table: "t", Name: "a", FieldName: "a", Type: MetricTypeValue},
				},
			},
			contains: "`t`.`a`",
		},
		{
			name:      "unknown type",
			metric:    &Metric{Name: "x", Type: MetricTypeUnknown},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.metric.Expression()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, got, tt.contains)
			}
		})
	}
}

func TestMetric_Alias(t *testing.T) {
	m := &Metric{Table: "orders", Name: "revenue", FieldName: "amount", Type: MetricTypeSum}
	alias, err := m.Alias()
	require.NoError(t, err)
	assert.Equal(t, "revenue", alias)
}

func TestMetric_Statement(t *testing.T) {
	m := &Metric{Table: "orders", Name: "cnt", FieldName: "*", Type: MetricTypeCount}
	stmt, err := m.Statement()
	require.NoError(t, err)
	assert.Contains(t, stmt, "AS cnt")
}

// ---- Join ----

func TestJoin_GetJoinType(t *testing.T) {
	t.Run("explicit join type", func(t *testing.T) {
		j := &Join{JoinType: "INNER JOIN"}
		assert.Equal(t, "INNER JOIN", j.GetJoinType())
	})

	t.Run("default left join", func(t *testing.T) {
		j := &Join{}
		assert.Equal(t, "LEFT JOIN", j.GetJoinType())
	})
}

// ---- DataSource (types package) ----

func TestTypesDataSource_Expression(t *testing.T) {
	t.Run("with database", func(t *testing.T) {
		ds := &DataSource{Database: "mydb", Name: "orders"}
		expr, err := ds.Expression()
		require.NoError(t, err)
		assert.Equal(t, "`mydb`.`orders`", expr)
	})

	t.Run("without database", func(t *testing.T) {
		ds := &DataSource{Name: "orders"}
		expr, err := ds.Expression()
		require.NoError(t, err)
		assert.Equal(t, "`orders`", expr)
	})
}

func TestTypesDataSource_Alias(t *testing.T) {
	t.Run("with alias", func(t *testing.T) {
		ds := &DataSource{Name: "orders", AliasName: "o"}
		alias, err := ds.Alias()
		require.NoError(t, err)
		assert.Equal(t, "o", alias)
	})

	t.Run("without alias defaults to name", func(t *testing.T) {
		ds := &DataSource{Name: "orders"}
		alias, err := ds.Alias()
		require.NoError(t, err)
		assert.Equal(t, "orders", alias)
	})
}

func TestTypesDataSource_Statement(t *testing.T) {
	ds := &DataSource{Database: "db", Name: "orders"}
	stmt, err := ds.Statement()
	require.NoError(t, err)
	assert.Equal(t, "`db`.`orders` AS orders", stmt)
}

// ---- Query / TimeInterval ----

func TestTimeInterval_ToFilter(t *testing.T) {
	ti := &TimeInterval{Name: "created_at", Start: "2026-01-01", End: "2026-02-01"}
	f := ti.ToFilter()
	assert.Equal(t, FilterOperatorTypeAnd, f.OperatorType)
	assert.Len(t, f.Children, 2)
	assert.Equal(t, FilterOperatorTypeGreaterEquals, f.Children[0].OperatorType)
	assert.Equal(t, FilterOperatorTypeLess, f.Children[1].OperatorType)
}

func TestQuery_TranslateTimeIntervalToFilter(t *testing.T) {
	t.Run("with valid time interval", func(t *testing.T) {
		q := &Query{
			TimeInterval: &TimeInterval{Name: "ts", Start: "2026-01-01", End: "2026-02-01"},
		}
		q.TranslateTimeIntervalToFilter()
		assert.Len(t, q.Filters, 1)
	})

	t.Run("with nil time interval", func(t *testing.T) {
		q := &Query{}
		q.TranslateTimeIntervalToFilter()
		assert.Empty(t, q.Filters)
	})

	t.Run("with empty start/end", func(t *testing.T) {
		q := &Query{
			TimeInterval: &TimeInterval{Name: "ts"},
		}
		q.TranslateTimeIntervalToFilter()
		assert.Empty(t, q.Filters)
	})
}

// ---- Result ----

func TestResult_SetDimensions(t *testing.T) {
	r := &Result{}
	q := &Query{
		Dimensions: []string{"city", "country"},
		Metrics:    []string{"revenue", "cnt"},
	}
	r.SetDimensions(q)
	assert.Equal(t, []string{"city", "country", "revenue", "cnt"}, r.Dimensions)
}

func TestResult_AddSource(t *testing.T) {
	r := &Result{}
	err := r.AddSource(map[string]any{"city": "beijing", "revenue": 100})
	require.NoError(t, err)
	assert.Len(t, r.Source, 1)
	assert.Equal(t, "beijing", r.Source[0]["city"])
}
