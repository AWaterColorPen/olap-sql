package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterOperatorType_IsTree(t *testing.T) {
	tests := []struct {
		op       FilterOperatorType
		expected bool
	}{
		{FilterOperatorTypeAnd, true},
		{FilterOperatorTypeOr, true},
		{FilterOperatorTypeEquals, false},
		{FilterOperatorTypeIn, false},
		{FilterOperatorTypeUnknown, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.op.IsTree(), "op=%v", tt.op)
	}
}

func TestFilter_Expression_SimpleOperators(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		expected string
	}{
		{
			name: "equals string",
			filter: &Filter{
				OperatorType: FilterOperatorTypeEquals,
				ValueType:    ValueTypeString,
				Name:         "city",
				Value:        []any{"beijing"},
			},
			expected: "city = 'beijing'",
		},
		{
			name: "equals integer",
			filter: &Filter{
				OperatorType: FilterOperatorTypeEquals,
				ValueType:    ValueTypeInteger,
				Name:         "age",
				Value:        []any{18},
			},
			expected: "age = 18",
		},
		{
			name: "in string list",
			filter: &Filter{
				OperatorType: FilterOperatorTypeIn,
				ValueType:    ValueTypeString,
				Name:         "status",
				Value:        []any{"active", "pending"},
			},
			expected: "status IN ('active', 'pending')",
		},
		{
			name: "not in",
			filter: &Filter{
				OperatorType: FilterOperatorTypeNotIn,
				ValueType:    ValueTypeString,
				Name:         "type",
				Value:        []any{"spam"},
			},
			expected: "type NOT IN ('spam')",
		},
		{
			name: "less equals",
			filter: &Filter{
				OperatorType: FilterOperatorTypeLessEquals,
				ValueType:    ValueTypeFloat,
				Name:         "score",
				Value:        []any{100.0},
			},
			expected: "score <= 100",
		},
		{
			name: "less",
			filter: &Filter{
				OperatorType: FilterOperatorTypeLess,
				ValueType:    ValueTypeInteger,
				Name:         "count",
				Value:        []any{50},
			},
			expected: "count < 50",
		},
		{
			name: "greater equals",
			filter: &Filter{
				OperatorType: FilterOperatorTypeGreaterEquals,
				ValueType:    ValueTypeInteger,
				Name:         "rank",
				Value:        []any{1},
			},
			expected: "rank >= 1",
		},
		{
			name: "greater",
			filter: &Filter{
				OperatorType: FilterOperatorTypeGreater,
				ValueType:    ValueTypeInteger,
				Name:         "id",
				Value:        []any{0},
			},
			expected: "id > 0",
		},
		{
			name: "like",
			filter: &Filter{
				OperatorType: FilterOperatorTypeLike,
				ValueType:    ValueTypeString,
				Name:         "name",
				Value:        []any{"%test%"},
			},
			expected: "name LIKE '%test%'",
		},
		{
			name: "has",
			filter: &Filter{
				OperatorType: FilterOperatorTypeHas,
				ValueType:    ValueTypeString,
				Name:         "tags",
				Value:        []any{"go"},
			},
			expected: "has(tags, 'go')",
		},
		{
			// FilterOperatorTypeExpression passes value through valueToStringSlice first,
			// so string values get quoted. Use ValueTypeUnknown with a numeric-like string
			// if you need a raw expression, or accept the quoting behavior.
			name: "expression with integer value",
			filter: &Filter{
				OperatorType: FilterOperatorTypeExpression,
				ValueType:    ValueTypeInteger,
				Name:         "x",
				Value:        []any{42},
			},
			expected: "42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.filter.Expression()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestFilter_Expression_TreeOperators(t *testing.T) {
	t.Run("and", func(t *testing.T) {
		f := &Filter{
			OperatorType: FilterOperatorTypeAnd,
			Children: []*Filter{
				{OperatorType: FilterOperatorTypeEquals, ValueType: ValueTypeString, Name: "city", Value: []any{"bj"}},
				{OperatorType: FilterOperatorTypeGreater, ValueType: ValueTypeInteger, Name: "age", Value: []any{18}},
			},
		}
		got, err := f.Expression()
		require.NoError(t, err)
		assert.Equal(t, "( city = 'bj' AND age > 18 )", got)
	})

	t.Run("or", func(t *testing.T) {
		f := &Filter{
			OperatorType: FilterOperatorTypeOr,
			Children: []*Filter{
				{OperatorType: FilterOperatorTypeEquals, ValueType: ValueTypeString, Name: "type", Value: []any{"A"}},
				{OperatorType: FilterOperatorTypeEquals, ValueType: ValueTypeString, Name: "type", Value: []any{"B"}},
			},
		}
		got, err := f.Expression()
		require.NoError(t, err)
		assert.Equal(t, "( type = 'A' OR type = 'B' )", got)
	})
}

func TestFilter_Expression_Errors(t *testing.T) {
	t.Run("unknown operator", func(t *testing.T) {
		f := &Filter{
			OperatorType: FilterOperatorTypeUnknown,
			ValueType:    ValueTypeString,
			Name:         "x",
			Value:        []any{"v"},
		}
		_, err := f.Expression()
		assert.Error(t, err)
	})

	t.Run("unsupported value type", func(t *testing.T) {
		f := &Filter{
			OperatorType: FilterOperatorTypeEquals,
			ValueType:    "UNSUPPORTED_TYPE",
			Name:         "x",
			Value:        []any{"v"},
		}
		_, err := f.Expression()
		assert.Error(t, err)
	})
}

func TestFilter_Statement(t *testing.T) {
	f := &Filter{
		OperatorType: FilterOperatorTypeEquals,
		ValueType:    ValueTypeString,
		Name:         "env",
		Value:        []any{"prod"},
	}
	got, err := f.Statement()
	require.NoError(t, err)
	assert.Equal(t, "env = 'prod'", got)
}

func TestFilter_Alias(t *testing.T) {
	f := &Filter{}
	_, err := f.Alias()
	assert.Error(t, err)
}

func TestTryToParseValue(t *testing.T) {
	assert.Equal(t, "'hello'", tryToParseValue("hello"))
	assert.Equal(t, "42", tryToParseValue(42))
	assert.Equal(t, "3.14", tryToParseValue(3.14))
	assert.Equal(t, "true", tryToParseValue(true)) // default branch
}

func TestFilter_ValueType_Unknown(t *testing.T) {
	// Unknown value type should auto-detect
	f := &Filter{
		OperatorType: FilterOperatorTypeEquals,
		ValueType:    ValueTypeUnknown,
		Name:         "x",
		Value:        []any{"hello"},
	}
	got, err := f.Expression()
	require.NoError(t, err)
	assert.Equal(t, "x = 'hello'", got)

	f2 := &Filter{
		OperatorType: FilterOperatorTypeEquals,
		ValueType:    "",
		Name:         "n",
		Value:        []any{99},
	}
	got2, err := f2.Expression()
	require.NoError(t, err)
	assert.Equal(t, "n = 99", got2)
}
