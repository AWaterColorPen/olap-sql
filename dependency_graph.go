package olapsql

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/models"
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

	current               string
	inDependencySourceKey map[string]bool
	isFold                map[string]bool
}

func (g *DependencyGraphBuilder) Build() (DependencyGraph, error) {
	g.isFold = g.buildFold()
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
	for _, child := range metric.Children {
		ok1 := g.isFold[child.Table]
		ok2 := metric.Table != child.Table
		if ok1 && ok2 {
			child.Children = nil
			child.Type = types.MetricTypeSum
			child.FieldName = child.Name
			child.Filter = nil
		}
	}
}

func (g *DependencyGraphBuilder) doIndependentDimension(dimension *types.Dimension) {
	for _, child := range dimension.Dependency {
		ok1 := g.isFold[child.Table]
		ok2 := dimension.Table != child.Table
		if ok1 && ok2 {
			child.Dependency = nil
			child.Type = types.DimensionTypeValue
		}
	}
}

func (g *DependencyGraphBuilder) buildFold() map[string]bool {
	graph, _ := GetDependencyTree(g.dictionary, g.current)
	return g.dfsFold(graph, g.current)
}

func (g *DependencyGraphBuilder) dfsFold(graph models.Graph, current string) map[string]bool {
	out := map[string]bool{current: false}
	source, _ := g.dictionary.GetSourceByKey(current)
	for _, node := range graph[current] {
		kv := g.dfsFold(graph, node)
		for k, v := range kv {
			out[k] = v
		}
		if source.Type == types.DataSourceTypeMergedJoin {
			out[node] = true
		}
	}
	for _, node := range graph[current] {
		if out[node] == true {
			out[current] = true
		}
	}
	return out
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
