package olapsql

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/models"
)

type Metrics []*models.Metric

func (m Metrics) NameIndex() map[string]*models.Metric {
	out := map[string]*models.Metric{}
	for _, v := range m {
		out[v.Name] = v
	}
	return out
}

func (m Metrics) KeyIndex() map[string]*models.Metric {
	out := map[string]*models.Metric{}
	for _, v := range m {
		out[v.GetKey()] = v
	}
	return out
}

type Dimensions []*models.Dimension

func (d Dimensions) NameIndex() map[string]*models.Dimension {
	out := map[string]*models.Dimension{}
	for _, v := range d {
		out[v.Name] = v
	}
	return out
}

func (d Dimensions) KeyIndex() map[string]*models.Dimension {
	out := map[string]*models.Dimension{}
	for _, v := range d {
		out[v.GetKey()] = v
	}
	return out
}

type joinNode struct {
	Children         []*joinNode
	metricNameMap    map[string]*models.Metric
	dimensionNameMap map[string]*models.Dimension
}

func (j *joinNode) FindMetric(name string) (*models.Metric, error) {
	m, ok := j.metricNameMap[name]
	if ok {
		return m, nil
	}

	for _, v := range j.Children {
		u, err := v.FindMetric(name)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		if m != nil {
			return nil, fmt.Errorf("duplicate metric name %v", name)
		}
		m = u
	}
	return m, nil
}

func (j *joinNode) FindDimension(name string) (*models.Dimension, error) {
	d, ok := j.dimensionNameMap[name]
	if ok {
		return d, nil
	}

	for _, v := range j.Children {
		u, err := v.FindDimension(name)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		if d != nil {
			return nil, fmt.Errorf("duplicate dimension name %v", name)
		}
		d = u
	}
	return d, nil
}

func newJoinNode(metrics []*models.Metric, dimensions []*models.Dimension) *joinNode {
	return &joinNode{
		metricNameMap:    Metrics(metrics).NameIndex(),
		dimensionNameMap: Dimensions(dimensions).NameIndex(),
	}
}

type joinTree struct {
	joinNode
	root string
}

func (j *joinTree) FindMetricByName(name string) (*models.Metric, error) {
	m, err := j.joinNode.FindMetric(name)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("not found metric name %v", name)
	}
	return m, nil
}

func (j *joinTree) FindDimensionByName(name string) (*models.Dimension, error) {
	d, err := j.joinNode.FindDimension(name)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, fmt.Errorf("not found dimension name %v", name)
	}
	return d, nil
}

type JoinTree interface {
	FindMetricByName(name string) (*models.Metric, error)
	FindDimensionByName(name string) (*models.Dimension, error)
}

type JoinTreeBuilder struct {
	tree       models.Graph
	root       string
	dictionary IAdapter
}

func (j *JoinTreeBuilder) Build() (JoinTree, error) {
	node, err := j.dfs(j.root)
	if err != nil {
		return nil, err
	}
	return &joinTree{joinNode: *node, root: j.root}, nil
}

func (j *JoinTreeBuilder) dfs(current string) (*joinNode, error) {
	metrics := j.dictionary.GetMetricsBySource(current)
	dimensions := j.dictionary.GetDimensionsBySource(current)

	node := newJoinNode(metrics, dimensions)
	for _, v := range j.tree[current] {
		child, err := j.dfs(v)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, child)
	}
	return node, nil
}
