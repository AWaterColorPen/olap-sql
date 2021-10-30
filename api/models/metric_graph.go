package models

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type metricGraph struct {
	name map[string]*types.Metric
}

func (m *metricGraph) GetByName(name string) (*types.Metric, error) {
	if v, ok := m.name[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("can't find %v in metric graph", name)
}

func (m *metricGraph) set(metric *types.Metric) {
	m.name[metric.Name] = metric
}

type MetricGraph interface {
	GetByName(name string) (*types.Metric, error)
}

type MetricGraphBuilder struct {
	dbType     types.DBType
	dictionary olapsql.IAdapter
	joinTree   JoinTree
}

func (m *MetricGraphBuilder) Build() (MetricGraph, error) {
	var key []string
	for _, v := range m.dictionary.GetMetric() {
		key = append(key, v.GetKey())
	}
	linq.From(key).Distinct().ToSlice(&key)

	queue, err := m.sort(key)
	if err != nil {
		return nil, err
	}

	graph := &metricGraph{name: map[string]*types.Metric{}}
	for _, k := range queue {
		metric, _ := m.dictionary.GetMetricByKey(k)
		current := &types.Metric{
			Type:      metric.Type,
			Table:     metric.DataSource,
			Name:      metric.Name,
			FieldName: metric.FieldName,
			Filter:    metric.Filter,
			DBType:    m.dbType,
		}
		for _, u := range metric.Composition {
			name := GetNameFromKey(u)
			value, _ := graph.GetByName(name)
			current.Children = append(current.Children, value)
		}
		graph.set(current)
	}

	return graph, nil
}

func (m *MetricGraphBuilder) sort(metrics []string) ([]string, error) {
	inDegree := map[string]int{}
	graph := map[string][]string{}

	for _, key := range metrics {
		name := GetNameFromKey(key)
		metric, err := m.joinTree.FindMetricByName(name)
		if err != nil {
			return nil, err
		}

		inDegree[metric.GetKey()] = 0
		for _, u := range metric.Composition {
			graph[u] = append(graph[u], metric.GetKey())
			inDegree[metric.GetKey()]++
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
	return queue, nil
}
