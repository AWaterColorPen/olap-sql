package models

import (
	"testing"

	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- Graph ----

func TestGraph_GetTree(t *testing.T) {
	t.Run("single node", func(t *testing.T) {
		g := Graph{"a": {"b", "c"}, "b": nil, "c": nil}
		tree := g.GetTree("a")
		assert.Equal(t, []string{"b", "c"}, tree["a"])
	})

	t.Run("multi level", func(t *testing.T) {
		g := Graph{"a": {"b"}, "b": {"c"}, "c": nil}
		tree := g.GetTree("a")
		assert.Equal(t, []string{"b"}, tree["a"])
		assert.Equal(t, []string{"c"}, tree["b"])
	})

	t.Run("empty start", func(t *testing.T) {
		g := Graph{}
		tree := g.GetTree("missing")
		assert.Empty(t, tree)
	})
}

// ---- DataSet ----

func TestDataSet_GetKey(t *testing.T) {
	ds := &DataSet{Name: "my_dataset"}
	assert.Equal(t, "my_dataset", ds.GetKey())
}

func TestDataSet_GetCurrent(t *testing.T) {
	ds := &DataSet{DataSource: "fact_orders"}
	assert.Equal(t, "fact_orders", ds.GetCurrent())
}

// ---- JoinPair ----

func TestJoinPair_GetJoinType(t *testing.T) {
	t.Run("get from first", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a", JoinType: "LEFT JOIN"},
			{DataSource: "b", JoinType: ""},
		}
		assert.Equal(t, "LEFT JOIN", pair.GetJoinType())
	})

	t.Run("get from second", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a", JoinType: ""},
			{DataSource: "b", JoinType: "INNER JOIN"},
		}
		assert.Equal(t, "INNER JOIN", pair.GetJoinType())
	})

	t.Run("both empty", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a"},
			{DataSource: "b"},
		}
		assert.Equal(t, "", pair.GetJoinType())
	})
}

func TestJoinPair_IsValid(t *testing.T) {
	t.Run("valid pair", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a", Dimension: []string{"id"}},
			{DataSource: "b", Dimension: []string{"user_id"}},
		}
		assert.NoError(t, pair.IsValid())
	})

	t.Run("wrong length", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a"},
		}
		assert.Error(t, pair.IsValid())
	})

	t.Run("mismatched dimensions", func(t *testing.T) {
		pair := JoinPair{
			{DataSource: "a", Dimension: []string{"id", "type"}},
			{DataSource: "b", Dimension: []string{"user_id"}},
		}
		assert.Error(t, pair.IsValid())
	})
}

// ---- DimensionJoins ----

func TestDimensionJoins_IsValid(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		dj := DimensionJoins{}
		assert.Error(t, dj.IsValid())
	})

	t.Run("valid", func(t *testing.T) {
		dj := DimensionJoins{
			{
				{DataSource: "fact", Dimension: []string{"user_id"}},
				{DataSource: "dim_user", Dimension: []string{"id"}},
			},
		}
		assert.NoError(t, dj.IsValid())
	})
}

func TestDimensionJoins_GetDependencyTree(t *testing.T) {
	t.Run("simple tree", func(t *testing.T) {
		dj := DimensionJoins{
			{
				{DataSource: "fact_orders", Dimension: []string{"user_id"}},
				{DataSource: "dim_users", Dimension: []string{"id"}},
			},
		}
		g, err := dj.GetDependencyTree("fact_orders")
		require.NoError(t, err)
		assert.NotEmpty(t, g)
	})
}

// ---- MergeJoin ----

func TestMergeJoin_IsValid(t *testing.T) {
	t.Run("too few", func(t *testing.T) {
		mj := MergeJoin{
			{DataSource: "fact"},
			{DataSource: "dim"},
		}
		assert.Error(t, mj.IsValid("fact"))
	})

	t.Run("wrong first source", func(t *testing.T) {
		mj := MergeJoin{
			{DataSource: "other", Dimension: []string{"id"}},
			{DataSource: "dim1", Dimension: []string{"id"}},
			{DataSource: "dim2", Dimension: []string{"id"}},
		}
		assert.Error(t, mj.IsValid("fact"))
	})

	t.Run("valid", func(t *testing.T) {
		mj := MergeJoin{
			{DataSource: "fact", Dimension: []string{"id"}},
			{DataSource: "dim1", Dimension: []string{"id"}},
			{DataSource: "dim2", Dimension: []string{"id"}},
		}
		assert.NoError(t, mj.IsValid("fact"))
	})

	t.Run("dimension mismatch", func(t *testing.T) {
		mj := MergeJoin{
			{DataSource: "fact", Dimension: []string{"id", "type"}},
			{DataSource: "dim1", Dimension: []string{"id"}},
			{DataSource: "dim2", Dimension: []string{"id"}},
		}
		assert.Error(t, mj.IsValid("fact"))
	})
}

func TestMergeJoin_GetDependencyTree(t *testing.T) {
	mj := MergeJoin{
		{DataSource: "fact"},
		{DataSource: "dim1"},
		{DataSource: "dim2"},
	}
	g, err := mj.GetDependencyTree("fact")
	require.NoError(t, err)
	assert.Contains(t, g["fact"], "dim1")
	assert.Contains(t, g["fact"], "dim2")
}

// ---- DataSource ----

func TestDataSource_GetKey(t *testing.T) {
	ds := &DataSource{Name: "fact_orders"}
	assert.Equal(t, "fact_orders", ds.GetKey())
}

func TestDataSource_IsFact(t *testing.T) {
	tests := []struct {
		dsType   types.DataSourceType
		expected bool
	}{
		{types.DataSourceTypeFact, true},
		{types.DataSourceTypeFactDimensionJoin, true},
		{types.DataSourceTypeMergeJoin, true},
		{types.DataSourceTypeDimension, false},
	}
	for _, tt := range tests {
		ds := &DataSource{Type: tt.dsType}
		assert.Equal(t, tt.expected, ds.IsFact(), "type=%v", tt.dsType)
	}
}

func TestDataSource_IsDimension(t *testing.T) {
	assert.True(t, (&DataSource{Type: types.DataSourceTypeDimension}).IsDimension())
	assert.False(t, (&DataSource{Type: types.DataSourceTypeFact}).IsDimension())
}

func TestDataSource_IsValid(t *testing.T) {
	t.Run("fact valid", func(t *testing.T) {
		ds := &DataSource{Type: types.DataSourceTypeFact, Name: "fact"}
		assert.NoError(t, ds.IsValid())
	})

	t.Run("dimension valid", func(t *testing.T) {
		ds := &DataSource{Type: types.DataSourceTypeDimension, Name: "dim"}
		assert.NoError(t, ds.IsValid())
	})

	t.Run("fact_dimension_join valid", func(t *testing.T) {
		ds := &DataSource{
			Name: "fact",
			Type: types.DataSourceTypeFactDimensionJoin,
			DimensionJoin: DimensionJoins{
				{
					{DataSource: "fact", Dimension: []string{"user_id"}},
					{DataSource: "dim_user", Dimension: []string{"id"}},
				},
			},
		}
		assert.NoError(t, ds.IsValid())
	})

	t.Run("merge_join valid", func(t *testing.T) {
		ds := &DataSource{
			Name: "fact",
			Type: types.DataSourceTypeMergeJoin,
			MergeJoin: MergeJoin{
				{DataSource: "fact", Dimension: []string{"id"}},
				{DataSource: "dim1", Dimension: []string{"id"}},
				{DataSource: "dim2", Dimension: []string{"id"}},
			},
		}
		assert.NoError(t, ds.IsValid())
	})

	t.Run("unknown type", func(t *testing.T) {
		ds := &DataSource{Type: "unknown", Name: "x"}
		assert.Error(t, ds.IsValid())
	})
}

func TestDataSource_GetDependencyTree(t *testing.T) {
	t.Run("fact", func(t *testing.T) {
		ds := &DataSource{Type: types.DataSourceTypeFact, Name: "fact"}
		g, err := ds.GetDependencyTree()
		require.NoError(t, err)
		assert.Contains(t, g, "fact")
	})

	t.Run("unknown type error", func(t *testing.T) {
		ds := &DataSource{Type: "unknown", Name: "x"}
		_, err := ds.GetDependencyTree()
		assert.Error(t, err)
	})
}

func TestDataSource_GetGetDependencyKey(t *testing.T) {
	t.Run("dimension join keys", func(t *testing.T) {
		ds := &DataSource{
			Name: "fact",
			Type: types.DataSourceTypeFactDimensionJoin,
			DimensionJoin: DimensionJoins{
				{
					{DataSource: "fact", Dimension: []string{"uid"}},
					{DataSource: "dim_user", Dimension: []string{"id"}},
				},
			},
		}
		keys := ds.GetGetDependencyKey()
		assert.Contains(t, keys, "fact")
		assert.Contains(t, keys, "dim_user")
	})

	t.Run("merge join keys", func(t *testing.T) {
		ds := &DataSource{
			Name: "fact",
			Type: types.DataSourceTypeMergeJoin,
			MergeJoin: MergeJoin{
				{DataSource: "fact"},
				{DataSource: "dim1"},
				{DataSource: "dim2"},
			},
		}
		keys := ds.GetGetDependencyKey()
		assert.Contains(t, keys, "dim1")
		assert.Contains(t, keys, "dim2")
	})
}

// ---- DataSources ----

func TestDataSources_KeyIndex(t *testing.T) {
	dss := DataSources{
		{Name: "fact_orders", Type: types.DataSourceTypeFact},
		{Name: "dim_users", Type: types.DataSourceTypeDimension},
	}
	idx := dss.KeyIndex()
	assert.Len(t, idx, 2)
	assert.Equal(t, "fact_orders", idx["fact_orders"].Name)
	assert.Equal(t, "dim_users", idx["dim_users"].Name)
}

// ---- Dimension ----

func TestDimension_GetKey(t *testing.T) {
	d := &Dimension{DataSource: "fact_orders", Name: "city"}
	assert.Equal(t, "fact_orders.city", d.GetKey())
}

func TestDimension_GetDependency(t *testing.T) {
	d := &Dimension{Dependency: []string{"a", "b"}}
	assert.Equal(t, []string{"a", "b"}, d.GetDependency())
}

// ---- Metric ----

func TestMetric_GetKey(t *testing.T) {
	m := &Metric{DataSource: "fact_orders", Name: "revenue"}
	assert.Equal(t, "fact_orders.revenue", m.GetKey())
}

func TestMetric_GetDependency(t *testing.T) {
	m := &Metric{Dependency: []string{"x", "y"}}
	assert.Equal(t, []string{"x", "y"}, m.GetDependency())
}

// ---- GetNameFromKey ----

func TestGetNameFromKey(t *testing.T) {
	assert.Equal(t, "city", GetNameFromKey("fact_orders.city"))
	assert.Equal(t, "revenue", GetNameFromKey("fact.revenue"))
	assert.Equal(t, "standalone", GetNameFromKey("standalone"))
}
