package models

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql"
)

type Metrics []*Metric

func (m Metrics) NameIndex() map[string]*Metric {
	out := map[string]*Metric{}
	for _, v := range m {
		out[v.Name] = v
	}
	return out
}

func (m Metrics) KeyIndex() map[string]*Metric {
	out := map[string]*Metric{}
	for _, v := range m {
		out[v.GetKey()] = v
	}
	return out
}

type Dimensions []*Dimension

func (d Dimensions) NameIndex() map[string]*Dimension {
	out := map[string]*Dimension{}
	for _, v := range d {
		out[v.Name] = v
	}
	return out
}

func (d Dimensions) KeyIndex() map[string]*Dimension {
	out := map[string]*Dimension{}
	for _, v := range d {
		out[v.GetKey()] = v
	}
	return out
}

type joinNode struct {
	Children         []*joinNode
	metricNameMap    map[string]*Metric
	dimensionNameMap map[string]*Dimension
}

func (j *joinNode) FindMetric(name string) (*Metric, error) {
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

func (j *joinNode) FindDimension(name string) (*Dimension, error) {
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

func newJoinNode(metrics []*Metric, dimensions []*Dimension) *joinNode {
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

func (j *joinTree) FindMetricByName(name string) (*Metric, error) {
	m, err := j.joinNode.FindMetric(name)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("not found metric name %v", name)
	}
	return m, nil
}

func (j *joinTree) FindDimensionByName(name string) (*Dimension, error) {
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
	Path(name string) ([]string, error)
	FindMetricByName(name string) (*Metric, error)
	FindDimensionByName(name string) (*Dimension, error)
}

type JoinTreeBuilder struct {
	tree       Graph
	root       string
	dictionary olapsql.IAdapter
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
