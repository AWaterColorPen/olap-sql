package models

import (
	"fmt"
	"strings"

	"github.com/awatercolorpen/olap-sql/api/types"
)

type IModel interface {
	GetKey() string
}

type DataSet struct {
	Name        string       `toml:"name"`
	DBType      types.DBType `toml:"type"`
	Description string       `toml:"description"`
	DataSource  string       `toml:"data_source"`
}

func (d *DataSet) GetKey() string {
	return d.Name
}

func (d *DataSet) GetCurrent() string {
	return d.DataSource
}

type Join struct {
	DataSource string   `toml:"data_source"`
	Dimension  []string `toml:"dimension"`
}

type JoinPair []*Join

func (j JoinPair) Get1() *Join {
	return j[0]
}

func (j JoinPair) Get2() *Join {
	return j[1]
}

func (j JoinPair) IsValid() error {
	if len(j) != 2 {
		return fmt.Errorf("join pair len %v != 2", len(j))
	}
	if len(j.Get1().Dimension) != len(j.Get2().Dimension) {
		return fmt.Errorf("join pair's dimension list len %v != %v", len(j.Get1().Dimension), len(j.Get2().Dimension))
	}
	return nil
}

type DimensionJoins []*JoinPair

func (d DimensionJoins) IsValid() error {
	if len(d) == 0 {
		return fmt.Errorf("dimension join len == 0")
	}
	for _, v := range d {
		if err := v.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

func (d DimensionJoins) GetDependencyTree(source string) (Graph, error) {
	g := Graph{source: nil}
	inDegree := map[string]int{}
	for _, v := range d {
		id1, id2 := v.Get1().DataSource, v.Get2().DataSource
		g[id1] = append(g[id1], id2)
		if _, ok := inDegree[id1]; !ok {
			inDegree[id1] = 0
		}
		inDegree[id2] = inDegree[id2] + 1
	}

	var root []string
	for k, v := range inDegree {
		if v > 1 {
			return nil, fmt.Errorf("it is graph, not a tree. id=%v has indegree=%v", k, v)
		}
		if v == 0 {
			root = append(root, k)
		}
	}
	if len(root) != 1 {
		return nil, fmt.Errorf("there are mulit root. id=%v", root)
	}

	g[source] = append(g[source], root...)
	return g, nil
}

type MergedJoin []*Join

func (m MergedJoin) IsValid(source string) error {
	if len(m) < 3 {
		return fmt.Errorf("merged join len %v < 3", len(m))
	}
	if m[0].DataSource != source {
		return fmt.Errorf("merged join's first dimension should be itself's, but %v", m[0].DataSource)
	}
	for i := 1; i < len(m); i++ {
		if len(m[i-1].Dimension) != len(m[i].Dimension) {
			return fmt.Errorf("merged join's dimension list len %v != %v", len(m[i-1].Dimension), len(m[i].Dimension))
		}
	}
	return nil
}

func (m MergedJoin) GetDependencyTree(source string) (Graph, error) {
	g := Graph{}
	for i := 1; i < len(m); i++ {
		g[source] = append(g[source], m[i].DataSource)
	}
	return g, nil
}

type DataSource struct {
	Database      string               `toml:"database"`
	Name          string               `toml:"name"`
	Alias         string               `toml:"alias"`
	Type          types.DataSourceType `toml:"type"`
	Description   string               `toml:"description"`
	DimensionJoin DimensionJoins       `toml:"dimension_join"`
	MergedJoin    MergedJoin           `toml:"merged_join"`
}

func (d *DataSource) GetKey() string {
	return d.Name
}

func (d *DataSource) IsFact() bool {
	switch d.Type {
	case types.DataSourceTypeFact, types.DataSourceTypeFactDimensionJoin, types.DataSourceTypeMergedJoin:
		return true
	default:
		return false
	}
}

func (d *DataSource) IsDimension() bool {
	switch d.Type {
	case types.DataSourceTypeDimension:
		return true
	default:
		return false
	}
}

func (d *DataSource) IsValid() error {
	switch d.Type {
	case types.DataSourceTypeFactDimensionJoin:
		return d.DimensionJoin.IsValid()
	case types.DataSourceTypeMergedJoin:
		return d.MergedJoin.IsValid(d.Name)
	case types.DataSourceTypeFact, types.DataSourceTypeDimension:
		return nil
	default:
		return fmt.Errorf("can't use datasource type=%v as dateset's datasource", d.Type)
	}
}

func (d *DataSource) GetDependencyTree() (Graph, error) {
	switch d.Type {
	case types.DataSourceTypeFactDimensionJoin:
		return d.DimensionJoin.GetDependencyTree(d.Name)
	case types.DataSourceTypeMergedJoin:
		return d.MergedJoin.GetDependencyTree(d.Name)
	case types.DataSourceTypeFact, types.DataSourceTypeDimension:
		return Graph{d.Name: nil}, nil
	default:
		return nil, fmt.Errorf("can't use datasource type=%v as dateset's datasource", d.Type)
	}
}

type DataSources []*DataSource

func (d DataSources) KeyIndex() map[string]*DataSource {
	out := map[string]*DataSource{}
	for _, v := range d {
		out[v.GetKey()] = v
	}
	return out
}

type Dimension struct {
	DataSource  string              `toml:"data_source"`
	Name        string              `toml:"name"`
	FieldName   string              `toml:"field_name"`
	Type        types.DimensionType `toml:"type"`
	ValueType   types.ValueType     `toml:"value_type"`
	Composition []string            `toml:"composition"`
	Description string              `toml:"description"`
}

func (d *Dimension) GetKey() string {
	return fmt.Sprintf("%v.%v", d.DataSource, d.Name)
}

type Metric struct {
	DataSource  string           `toml:"data_source"`
	Name        string           `toml:"name"`
	FieldName   string           `toml:"field_name"`
	Type        types.MetricType `toml:"type"`
	ValueType   types.ValueType  `toml:"value_type"`
	Description string           `toml:"description"`
	Composition []string         `toml:"composition"`
	Filter      *types.Filter    `toml:"filter"`
}

func (m *Metric) GetKey() string {
	return fmt.Sprintf("%v.%v", m.DataSource, m.Name)
}

func GetNameFromKey(key string) string {
	out := strings.Split(key, ".")
	return out[len(out)-1]
}
