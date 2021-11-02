package olapsql

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type dependencyGraph struct {
	key map[string]interface{}
}

func (g *dependencyGraph) Get(key string) (interface{}, error) {
	if v, ok := g.key[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("can't find %v in dependency graph", key)
}

func (g *dependencyGraph) GetMetric(key string) (*types.Metric, error) {
	v, err := g.Get(key)
	if err != nil {
		return nil, err
	}
	if u, ok := v.(*types.Metric); ok {
		return u, nil
	}
	return nil, fmt.Errorf("found item, but it is not a metric")
}

func (g *dependencyGraph) GetDimension(key string) (*types.Dimension, error) {
	v, err := g.Get(key)
	if err != nil {
		return nil, err
	}
	if u, ok := v.(*types.Dimension); ok {
		return u, nil
	}
	return nil, fmt.Errorf("found item, but it is not a dimension")
}

func (g *dependencyGraph) set(in interface{}) {
	switch v := in.(type) {
	case *types.Metric:
		k := fmt.Sprintf("%v.%v", v.Table, v.Name)
		g.key[k] = in
	case *types.Dimension:
		k := fmt.Sprintf("%v.%v", v.Table, v.Name)
		g.key[k] = in
	}
}

type DependencyGraph interface {
	Get(key string) (interface{}, error)
	GetMetric(key string) (*types.Metric, error)
	GetDimension(key string) (*types.Dimension, error)
}

type DependencyGraphBuilder struct {
	dbType     types.DBType
	dictionary IAdapter

	current    string
	inDependencySourceKey map[string]bool
}

func (g *DependencyGraphBuilder) Build() (DependencyGraph, error) {
	g.inDependencySourceKey = map[string]bool{g.current: true}
	source, _ := g.dictionary.GetSourceByKey(g.current)
	for _, in := range source.GetGetDependencyKey() {
		g.inDependencySourceKey[in] = true
	}
	graph := &dependencyGraph{key: map[string]interface{}{}}
	if err := g.buildMetricDependency(graph); err != nil {
		return nil, err
	}
	if err := g.buildDimensionDependency(graph); err != nil {
		return nil, err
	}
	return graph, nil
}

func (g *DependencyGraphBuilder) buildMetricDependency(graph *dependencyGraph) error {
	var key []string
	for _, v := range g.dictionary.GetMetric() {
		key = append(key, v.GetKey())
	}

	key2Model := func(k string) iDependencyModel {
		m, _ := g.dictionary.GetMetricByKey(k)
		return m
	}
	queue := g.sort(key, key2Model)
	for _, k := range queue {
		metric, _ := g.dictionary.GetMetricByKey(k)
		current := &types.Metric{
			Table:     metric.DataSource,
			Name:      metric.Name,
			FieldName: metric.FieldName,
			Type:      metric.Type,
			ValueType: metric.ValueType,
			Filter:    metric.Filter,
			DBType:    g.dbType,
		}
		for _, u := range metric.Composition {
			value, _ := graph.GetMetric(u)
			current.Children = append(current.Children, value)
		}
		g.doIndependentMetric(current)
		graph.set(current)
	}
	return nil
}

func (g *DependencyGraphBuilder) buildDimensionDependency(graph *dependencyGraph) error {
	var key []string
	for _, v := range g.dictionary.GetDimension() {
		key = append(key, v.GetKey())
	}

	key2Model := func(k string) iDependencyModel {
		m, _ := g.dictionary.GetDimensionByKey(k)
		return m
	}
	queue := g.sort(key, key2Model)
	for _, k := range queue {
		dimension, _ := g.dictionary.GetDimensionByKey(k)
		current := &types.Dimension{
			Table:     dimension.DataSource,
			Name:      dimension.Name,
			Type:      dimension.Type,
			ValueType: dimension.ValueType,
			FieldName: dimension.FieldName,
		}
		for _, u := range dimension.Composition {
			value, _ := graph.GetDimension(u)
			current.Dependency = append(current.Dependency, value)
		}
		g.doIndependentDimension(current)
		graph.set(current)
	}
	return nil
}

func (g *DependencyGraphBuilder) doIndependentMetric(metric *types.Metric) {
	independent := true
	for _, child := range metric.Children {
		_, ok := g.inDependencySourceKey[child.Table]
		independent = independent && ok
	}
	if independent {
		return
	}
	metric.Children = nil
	metric.Type = types.MetricTypeValue
}

func (g *DependencyGraphBuilder) doIndependentDimension(dimension *types.Dimension) {
	independent := true
	for _, child := range dimension.Dependency {
		_, ok := g.inDependencySourceKey[child.Table]
		independent = independent && ok
	}
	if independent {
		return
	}
	dimension.Dependency = nil
	dimension.Type = types.DimensionTypeValue
}

type iDependencyModel interface {
	GetKey() string
	GetDependency() []string
}

func (g *DependencyGraphBuilder) sort(key []string, key2Model func(string) iDependencyModel) []string {
	inDegree := map[string]int{}
	graph := map[string][]string{}

	for _, k := range key {
		iModel := key2Model(k)
		inDegree[iModel.GetKey()] = 0
		for _, u := range iModel.GetDependency() {
			graph[u] = append(graph[u], iModel.GetKey())
			inDegree[iModel.GetKey()]++
			if _, ok := inDegree[u]; !ok {
				inDegree[u] = 0
			}
		}
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
	return queue
}
