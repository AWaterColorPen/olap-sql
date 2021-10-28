package models

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type IModel interface {
	GetKey() string
}

type DataSet struct {
	Name        string           `toml:"name"`
	DBType      types.DBType     `toml:"type"`
	Description string           `toml:"description"`
	Root        string           `toml:"root"`
	Join        []*DataSetJoin   `toml:"join"`
	Merged      []*DataSetMerged `toml:"merged"`
}

func (d *DataSet) GetKey() string {
	return d.Name
}

func (d *DataSet) GetDataSource() []string {
	var out []string
	for _, join := range d.Join {
	    out = append(out, join.DataSource1, join.DataSource2)
	}
	linq.From(out).Distinct().ToSlice(&out)
	return out
}

func (d *DataSet) GetRoot() string {
	return d.Root
}

func (d *DataSet) JoinTopologyGraph() (Graph, error) {
	inDegree := map[string]int{}
	graph := Graph{}

	for _, v := range d.Join {
		if _, ok := inDegree[v.DataSource1]; !ok {
			inDegree[v.DataSource1] = 0
		}
		inDegree[v.DataSource2] = inDegree[v.DataSource2] + 1
		graph[v.DataSource1] = append(graph[v.DataSource1], v.DataSource2)
	}

	var queue []string
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}
	for i := 0; i < len(queue); i++ {
		node := queue[i]
		for _, v := range graph[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(inDegree) != len(queue) {
		return nil, fmt.Errorf("it is not a topology graph. node=%v, intop=%v", len(inDegree), len(queue))
	}
	return graph, nil
}

type DataSetMerged struct {
}

type DataSetJoin struct {
	DataSource1 string    `toml:"data_source1"`
	DataSource2 string    `toml:"data_source2"`
	JoinOn      []*JoinOn `toml:"join_on"`
}

type JoinOn struct {
	Dimension1 string `toml:"dimension1"`
	Dimension2 string `toml:"dimension2"`
}

type DataSource struct {
	Database    string `toml:"database"`
	Name        string `toml:"name"`
	Alias       string `toml:"alias"`
	Description string `toml:"description"`
}

func (d *DataSource) GetKey() string {
	return d.Name
}

type Dimension struct {
	DataSource  string          `toml:"data_source"`
	Name        string          `toml:"name"`
	FieldName   string          `toml:"field_name"`
	ValueType   types.ValueType `toml:"value_type"`
	Description string          `toml:"description"`
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
	Composition []string         `toml:"composition"`
	Description string           `toml:"description"`
	Filter      *types.Filter    `toml:"filter"`
}

func (m *Metric) GetKey() string {
	return fmt.Sprintf("%v.%v", m.DataSource, m.Name)
}

func GetNameFromKey(key string) string {
	out := strings.Split(key, ".")
	return out[len(out)-1]
}
