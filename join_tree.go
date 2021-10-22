package olapsql

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
)

type Metrics []*models.Metric

func (m Metrics) IdIndex() map[uint64]*models.Metric {
	out := map[uint64]*models.Metric{}
	for _, v := range m {
		out[v.ID] = v
	}
	return out
}

func (m Metrics) NameIndex() map[string]*models.Metric {
	out := map[string]*models.Metric{}
	for _, v := range m {
		out[v.Name] = v
	}
	return out
}

type Dimensions []*models.Dimension

func (d Dimensions) IdIndex() map[uint64]*models.Dimension {
	out := map[uint64]*models.Dimension{}
	for _, v := range d {
		out[v.ID] = v
	}
	return out
}

func (d Dimensions) NameIndex() map[string]*models.Dimension {
	out := map[string]*models.Dimension{}
	for _, v := range d {
		out[v.Name] = v
	}
	return out
}

type joinNode struct {
	Children         []*joinNode
	source           *models.DataSource
	metricNameMap    map[string]*models.Metric
	dimensionNameMap map[string]*models.Dimension
}

func (j *joinNode) ID() uint64 {
	return j.source.ID
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

func newJoinNode(metrics []*models.Metric, dimensions []*models.Dimension, source *models.DataSource) *joinNode {
	return &joinNode{
		source:           source,
		metricNameMap:    Metrics(metrics).NameIndex(),
		dimensionNameMap: Dimensions(dimensions).NameIndex(),
	}
}

type joinTree struct {
	joinNode
	root     uint64
	inverted map[uint64]uint64
}

func (j *joinTree) Path(current uint64) ([]uint64, error) {
	var out []uint64
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

func (j *joinTree) FindMetric(name string) (*models.Metric, error) {
	m, err := j.joinNode.FindMetric(name)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, fmt.Errorf("not found metric name %v", name)
	}
	return m, nil
}

func (j *joinTree) FindDimension(name string) (*models.Dimension, error) {
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
	Path(uint64) ([]uint64, error)
	FindMetric(string) (*models.Metric, error)
	FindDimension(string) (*models.Dimension, error)
}

type JoinTreeBuilder struct {
	tree      map[uint64][]uint64
	root      uint64
	metrics []*models.Metric
	dimensions []*models.Dimension
	sourceMap map[uint64]*models.DataSource
}

func (j *JoinTreeBuilder) Build() (JoinTree, error) {
	node, err := j.dfs(j.root)
	if err != nil {
		return nil, err
	}

	inverted := map[uint64]uint64{}
	for k, v := range j.tree {
		for _, u := range v {
			inverted[u] = k
		}
	}
	return &joinTree{joinNode: *node, root: j.root, inverted: inverted}, nil
}

func (j *JoinTreeBuilder) dfs(current uint64) (*joinNode, error) {
	source, ok := j.sourceMap[current]
	if !ok {
		return nil, fmt.Errorf("can't find %v in source map", current)
	}
	metrics, dimensions := j.GetSourceRelateMetricsAndDimensions(source)
	node := newJoinNode(metrics, dimensions, source)
	for _, v := range j.tree[current] {
		child, err := j.dfs(v)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, child)
	}
	return node, nil
}

// GetSourceRelateMetricsAndDimensions 获取 Source 的 Metrics 和 Dimensions
func (j *JoinTreeBuilder) GetSourceRelateMetricsAndDimensions(src *models.DataSource) ([]*models.Metric, []*models.Dimension) {
	metrics := make([]*models.Metric, 0)
	dimensions := make([]*models.Dimension, 0)
	for _, metric := range j.metrics {
		if metric.DataSourceID == src.ID {
			metrics = append(metrics, metric)
		}
	}
	for _, dimension := range j.dimensions {
		if dimension.DataSourceID == src.ID {
			dimensions = append(dimensions, dimension)
		}
	}
	return metrics, dimensions
}
