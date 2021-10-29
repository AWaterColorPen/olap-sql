package olapsql

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
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

func (j *joinNode) FindMetric(key string) (*models.Metric, error) {
	m, ok := j.metricNameMap[key]
	if ok {
		return m, nil
	}

	for _, v := range j.Children {
		u, err := v.FindMetric(key)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		if m != nil {
			return nil, fmt.Errorf("duplicate metric key %v", key)
		}
		m = u
	}
	return m, nil
}

func (j *joinNode) FindDimension(key string) (*models.Dimension, error) {
	d, ok := j.dimensionNameMap[key]
	if ok {
		return d, nil
	}

	for _, v := range j.Children {
		u, err := v.FindDimension(key)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		if d != nil {
			return nil, fmt.Errorf("duplicate dimension name %v", key)
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
	root     string
	inverted map[string]string
}

func (j *joinTree) Path(current string) ([]string, error) {
	var out []string
	for true {
		out = append(out, current)
		if current == j.root {
			break
		}
		u, ok := j.inverted[current]
		if !ok {
			return nil, fmt.Errorf("can't find %v node", current)
		}
		current = u
	}
	linq.From(out).Reverse().ToSlice(&out)
	return out, nil
}

func (j *joinTree) FindMetricByName(key string) (*models.Metric, error) {
	m, err := j.joinNode.FindMetric(key)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("not found metric name %v", key)
	}
	return m, nil
}

func (j *joinTree) FindDimensionByName(key string) (*models.Dimension, error) {
	d, err := j.joinNode.FindDimension(key)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, fmt.Errorf("not found dimension key %v", key)
	}
	return d, nil
}

type JoinTree interface {
	Path(key string) ([]string, error)
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

	inverted := map[string]string{}
	for k, v := range j.tree {
		for _, u := range v {
			inverted[u] = k
		}
	}
	return &joinTree{joinNode: *node, root: j.root, inverted: inverted}, nil
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
