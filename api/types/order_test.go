package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderBy_Expression(t *testing.T) {
	tests := []struct {
		name      string
		order     *OrderBy
		expected  string
		expectErr bool
	}{
		{
			name: "dimension ascending",
			order: &OrderBy{
				Table:         "users",
				Name:          "created_at",
				FieldProperty: FieldPropertyDimension,
				Direction:     OrderDirectionTypeAscending,
			},
			expected: "created_at ASC",
		},
		{
			name: "metric descending",
			order: &OrderBy{
				Table:         "orders",
				Name:          "total_amount",
				FieldProperty: FieldPropertyMetric,
				Direction:     OrderDirectionTypeDescending,
			},
			expected: "total_amount DESC",
		},
		{
			name: "unknown direction",
			order: &OrderBy{
				Name:          "x",
				FieldProperty: FieldPropertyDimension,
				Direction:     OrderDirectionTypeUnknown,
			},
			expectErr: true,
		},
		{
			name: "unsupported field property",
			order: &OrderBy{
				Name:          "x",
				FieldProperty: "UNKNOWN",
				Direction:     OrderDirectionTypeAscending,
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.order.Expression()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestOrderBy_Statement(t *testing.T) {
	o := &OrderBy{
		Name:          "score",
		FieldProperty: FieldPropertyMetric,
		Direction:     OrderDirectionTypeDescending,
	}
	got, err := o.Statement()
	require.NoError(t, err)
	assert.Equal(t, "score DESC", got)
}

func TestOrderBy_Alias(t *testing.T) {
	o := &OrderBy{}
	_, err := o.Alias()
	assert.Error(t, err)
}
