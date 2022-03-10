package olapsql

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/types"
	"sort"
)

type columnStruct struct {
	FieldProperty types.FieldProperty
	Metric        *types.Metric
	Dimension     *types.Dimension
}

func (c *columnStruct) GetTable() string {
	if c.FieldProperty == types.FieldPropertyMetric {
		return c.Metric.Table
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		return c.Dimension.Table
	}
	return ""
}

func (c *columnStruct) GetName() string {
	if c.FieldProperty == types.FieldPropertyMetric {
		return c.Metric.Name
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		return c.Dimension.Name
	}
	return ""
}

func (c *columnStruct) GetValueType() types.ValueType {
	if c.FieldProperty == types.FieldPropertyMetric {
		return c.Metric.ValueType
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		return c.Dimension.ValueType
	}
	return ""
}

func (c *columnStruct) GetExpression() string {
	if c.FieldProperty == types.FieldPropertyMetric {
		expression, _ := c.Metric.Expression()
		return expression
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		expression, _ := c.Dimension.Expression()
		return expression
	}
	return ""
}

func (c *columnStruct) IsAs() bool {
	if c.FieldProperty == types.FieldPropertyMetric {
		return c.Metric.Type == types.MetricTypeAs
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		return c.Dimension.Type == types.DimensionTypeMulti || c.Dimension.Type == types.DimensionTypeCase
	}
	return false
}

func (c *columnStruct) GetAsTables() []string {
	var out []string
	if c.FieldProperty == types.FieldPropertyMetric {
		out = append(out, c.Metric.Children[0].Table)
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		for _, v := range c.Dimension.Dependency {
			out = append(out, v.Table)
		}
	}
	linq.From(out).Distinct().ToSlice(&out)
	return out
}

func (c *columnStruct) GetTables() []string {
	if c.IsAs() {
		return c.GetAsTables()
	}
	return []string{c.GetTable()}
}

func (c *columnStruct) GetTargetFieldName(target string) string {
	if c.FieldProperty == types.FieldPropertyMetric {
		for _, v := range c.Metric.Children {
			if v.Table == target {
				return v.Name
			}
		}
	}
	if c.FieldProperty == types.FieldPropertyDimension {
		for _, v := range c.Dimension.Dependency {
			if v.Table == target {
				return v.Name
			}
		}
	}
	return ""
}

func isSameColumnTables(t1, t2 []string) error {
	if len(t1) != len(t2) {
		return fmt.Errorf("table is not same t1=%v, t2=%v", t1, t2)
	}
	sort.Strings(t1)
	sort.Strings(t2)
	// TODO linq
	for i := 0; i < len(t1); i++ {
		if t1[i] != t2[i] {
			return fmt.Errorf("table is not same t1=%v, t2=%v", t1, t2)
		}
	}
	return nil
}
