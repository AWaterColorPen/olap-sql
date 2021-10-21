package dictionary

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type metricGraph struct {
	id   map[uint64]*types.Metric
	name map[string]*types.Metric
}

func (m *metricGraph) GetByID(id uint64) (*types.Metric, error) {
	if v, ok := m.id[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("can't find %v in metric graph", id)
}

func (m *metricGraph) GetByName(name string) (*types.Metric, error) {
	if v, ok := m.name[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("can't find %v in metric graph", name)
}

func (m *metricGraph) set(id uint64, metric *types.Metric) {
	m.id[id] = metric
	m.name[metric.Name] = metric
}

type MetricGraph interface {
	GetByID(id uint64) (*types.Metric, error)
	GetByName(name string) (*types.Metric, error)
}

type MetricGraphBuilder struct {
	sourceMap map[uint64]*models.DataSource
	metricMap map[uint64]*models.Metric
	joinTree  JoinTree
}

func (m *MetricGraphBuilder) Build() (MetricGraph, error) {
	var name []string
	for _, v := range m.metricMap {
		name = append(name, v.Name)
	}
	linq.From(name).Distinct().ToSlice(&name)

	queue, err := m.sort(name)
	if err != nil {
		return nil, err
	}

	graph := &metricGraph{id: map[uint64]*types.Metric{}, name: map[string]*types.Metric{}}
	for _, v := range queue {
		metric := m.metricMap[v]
		source := m.sourceMap[metric.DataSourceID]
		current := &types.Metric{
			Type:      metric.Type,
			Table:     source.GetTableName(),
			Name:      metric.Name,
			FieldName: metric.FieldName,
			If:        metric.If,
			DBType:    source.Type,
		}
		if metric.Composition != nil {
			for _, u := range metric.Composition.MetricID {
				value, _ := graph.GetByID(u)
				current.Children = append(current.Children, value)
			}
		}
		graph.set(v, current)
	}

	return graph, nil
}

func (m *MetricGraphBuilder) sort(metrics []string) ([]uint64, error) {
	inDegree := map[uint64]int{}
	graph := map[uint64][]uint64{}

	for _, v := range metrics {
		metric, err := m.joinTree.FindMetric(v)
		if err != nil {
			return nil, err
		}

		inDegree[metric.ID] = 0
		if metric.Composition != nil {
			for _, u := range metric.Composition.MetricID {
				graph[u] = append(graph[u], metric.ID)
				inDegree[metric.ID]++
				if _, ok := inDegree[u]; !ok {
					inDegree[u] = 0
				}
			}
		}
	}

	var queue []uint64
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
