package olapsql

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/awatercolorpen/olap-sql/api/models"
)

type AdapterType string

const (
	FILEAdapter AdapterType = "FILE"
)

// IAdapter Adapter适配器
type IAdapter interface {
	BuildDataSourceAdapter(string) (IAdapter, error)

	GetMetric() []*models.Metric

	GetDataSetByKey(string) (*models.DataSet, error)
	GetSourceByKey(string) (*models.DataSource, error)
	GetMetricByKey(string) (*models.Metric, error)
	GetDimensionByKey(string) (*models.Dimension, error)

	GetMetricsBySource(string) []*models.Metric
	GetDimensionsBySource(string) []*models.Dimension
}

func IsValidJoin(adapter IAdapter, join *models.Join) error {
	for _, on := range join.Dimension {
		k := fmt.Sprintf("%v.%v", join.DataSource, on)
		if _, err := adapter.GetDimensionByKey(k); err != nil {
			return fmt.Errorf("invalid dataset join setting=%v, err=%v", join.DataSource, err)
		}
	}
	return nil
}

func GetDependencyTree(adapter IAdapter, current string) (models.Graph, error) {
	if root, err := adapter.GetSourceByKey(current); err != nil {
		return nil, err
	} else if root.IsDimension() {
		return nil, fmt.Errorf("can't get dependency from a dimension datasource")
	}

	graph := models.Graph{current: nil}
	queue := []string{current}
	for i := 0; i < len(queue); i++ {
		node, err := adapter.GetSourceByKey(queue[i])
		if err != nil {
			return nil, err
		}
		tree, err := node.GetDependencyTree()
		if err != nil {
			return nil, err
		}
		for k, v := range tree {
			for _, u := range v {
				queue = append(queue, u)
				graph[k] = append(graph[k], u)
			}
		}
	}
	//
	// for _, v := range d.DimensionJoin {
	// 	if _, ok := inDegree[v.DataSource1]; !ok {
	// 		inDegree[v.DataSource1] = 0
	// 	}
	// 	inDegree[v.DataSource2] = inDegree[v.DataSource2] + 1
	// 	graph[v.DataSource1] = append(graph[v.DataSource1], v.DataSource2)
	// }
	//
	// for i := 0; i < len(queue); i++ {
	// 	node := queue[i]
	// 	for _, v := range graph[node] {
	// 		inDegree[v]--
	// 		if inDegree[v] == 0 {
	// 			queue = append(queue, v)
	// 		}
	// 	}
	// }
	// if len(inDegree) != len(queue) {
	// 	return nil, fmt.Errorf("it is not a topology graph. node=%v, intop=%v", len(inDegree), len(queue))
	// }
	return graph, nil
}

// AdapterOption Adapter配置
type AdapterOption struct {
	Type AdapterType `json:"type"`
	Dsn  string      `json:"dsn"`
}

func NewAdapter(option *AdapterOption) (IAdapter, error) {
	switch option.Type {
	case FILEAdapter:
		return newDictionaryAdapterByFile(option)
	default:
		return nil, fmt.Errorf("not supported adapter type %v", option.Type)
	}
}

// FileAdapter 文件适配器
type FileAdapter struct {
	Sets       []*models.DataSet
	Sources    []*models.DataSource
	Metrics    []*models.Metric
	Dimensions []*models.Dimension
}

func (f *FileAdapter) BuildDataSourceAdapter(key string) (IAdapter, error) {
	source, err := f.GetSourceByKey(key)
	if err != nil {
		return nil, err
	}
	tree, err := source.GetDependencyTree()
	if err != nil {
		return nil, err
	}

	var sKey []string
	for k, v := range tree {
		sKey = append(sKey, k)
		sKey = append(sKey, v...)
	}

	sources := f.getSourcesByKeys(sKey)
	metrics := f.getMetricsBySourceKeys(sKey)
	dimensions := f.getDimensionsBySourceKeys(sKey)

	adapter := &FileAdapter{
		Sources:    sources,
		Metrics:    metrics,
		Dimensions: dimensions,
	}
	return adapter, adapter.isValid()
}

func (f *FileAdapter) GetMetric() []*models.Metric {
	return f.Metrics
}

func (f *FileAdapter) GetDataSetByKey(key string) (*models.DataSet, error) {
	for _, set := range f.Sets {
		if set.GetKey() == key {
			return set, nil
		}
	}
	return nil, fmt.Errorf("can not find '%v' data set", key)
}

func (f *FileAdapter) GetSourceByKey(key string) (*models.DataSource, error) {
	for _, source := range f.Sources {
		if source.GetKey() == key || source.Alias == key {
			return source, nil
		}
	}
	return nil, fmt.Errorf("can not find '%v' data source", key)
}

func (f *FileAdapter) GetMetricByKey(key string) (*models.Metric, error) {
	for _, metric := range f.Metrics {
		if metric.GetKey() == key {
			return metric, nil
		}
	}
	return nil, fmt.Errorf("can not find '%v' metric", key)
}

func (f *FileAdapter) GetDimensionByKey(key string) (*models.Dimension, error) {
	for _, dimension := range f.Dimensions {
		if dimension.GetKey() == key {
			return dimension, nil
		}
	}
	return nil, fmt.Errorf("can not find '%v' dimension", key)
}

func (f *FileAdapter) GetMetricsBySource(key string) []*models.Metric {
	var out []*models.Metric
	for _, metric := range f.Metrics {
		if metric.DataSource == key {
			out = append(out, metric)
		}
	}
	return out
}

func (f *FileAdapter) GetDimensionsBySource(key string) []*models.Dimension {
	var out []*models.Dimension
	for _, dimension := range f.Dimensions {
		if dimension.DataSource == key {
			out = append(out, dimension)
		}
	}
	return out
}

func (f *FileAdapter) getSourcesByKeys(key []string) []*models.DataSource {
	set := getKeySet(key)
	var out []*models.DataSource
	for _, source := range f.Sources {
		if _, ok := set[source.GetKey()]; ok {
			out = append(out, source)
		} else if _, ok := set[source.Alias]; ok {
			out = append(out, source)
		}
	}
	return out
}

func (f *FileAdapter) getMetricsByKeys(key []string) []*models.Metric {
	set := getKeySet(key)
	var out []*models.Metric
	for _, metric := range f.Metrics {
		if _, ok := set[metric.GetKey()]; ok {
			out = append(out, metric)
		}
	}
	return out
}

func (f *FileAdapter) getMetricsBySourceKeys(key []string) []*models.Metric {
	set := getKeySet(key)
	var out []*models.Metric
	for _, metric := range f.Metrics {
		if _, ok := set[metric.DataSource]; ok {
			out = append(out, metric)
		}
	}
	return out
}

func (f *FileAdapter) getDimensionsByKeys(key []string) []*models.Dimension {
	set := getKeySet(key)
	var out []*models.Dimension
	for _, dimension := range f.Dimensions {
		if _, ok := set[dimension.GetKey()]; ok {
			out = append(out, dimension)
		}
	}
	return out
}

func (f *FileAdapter) getDimensionsBySourceKeys(key []string) []*models.Dimension {
	set := getKeySet(key)
	var out []*models.Dimension
	for _, dimension := range f.Dimensions {
		if _, ok := set[dimension.DataSource]; ok {
			out = append(out, dimension)
		}
	}
	return out
}

func (f *FileAdapter) isValidDataSet(set *models.DataSet) error {
	root, err := f.GetSourceByKey(set.GetCurrent())
	if err != nil {
		return err
	}

	if !root.IsFact() {
		return fmt.Errorf("can't use one datasource with type=%v as dateset's datasource", root.Type)
	}

	return nil
}

func (f *FileAdapter) isValidDataSource(source *models.DataSource) error {
	if err := source.IsValid(); err != nil {
		return err
	}
	for _, v := range source.DimensionJoin {
		for _, u := range []*models.Join{v.Get1(), v.Get2()} {
			if err := IsValidJoin(f, u); err != nil {
				return err
			}
		}
	}
	for _, v := range source.MergedJoin {
		if err := IsValidJoin(f, v); err != nil {
			return err
		}
	}

	if !source.IsDimension() {
		if _, err := GetDependencyTree(f, source.GetKey()); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileAdapter) isValidMetric(metric *models.Metric) error {
	if _, err := f.GetSourceByKey(metric.DataSource); err != nil {
		return err
	}
	for _, composition := range metric.Composition {
		if _, err := f.GetMetricByKey(composition); err != nil {
			return err
		}
	}
	// if metric.Filter != nil {
	// 	if _, err := f.GetMetricByKey(metric.Filter.Name); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (f *FileAdapter) isValidDimension(dimension *models.Dimension) error {
	_, err := f.GetSourceByKey(dimension.DataSource)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileAdapter) isValid() error {
	for _, set := range f.Sets {
		if err := f.isValidDataSet(set); err != nil {
			return err
		}
	}
	for _, source := range f.Sources {
		if err := f.isValidDataSource(source); err != nil {
			return err
		}
	}
	for _, metric := range f.Metrics {
		if err := f.isValidMetric(metric); err != nil {
			return err
		}
	}
	for _, dimensions := range f.Dimensions {
		if err := f.isValidDimension(dimensions); err != nil {
			return err
		}
	}
	return nil
}

func newDictionaryAdapterByFile(option *AdapterOption) (*FileAdapter, error) {
	b, err := ioutil.ReadFile(option.Dsn)
	if err != nil {
		return nil, err
	}

	adapter := &FileAdapter{}
	extension := filepath.Ext(option.Dsn)
	switch extension {
	case ".toml":
		if err = toml.Unmarshal(b, adapter); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("not supported extension %v", extension)
	}

	if err = adapter.isValid(); err != nil {
		return nil, err
	}
	return adapter, nil
}

func getKeySet(key []string) map[interface{}]struct{} {
	set := make(map[interface{}]struct{})
	for _, k := range key {
		set[k] = struct{}{}
	}
	return set
}
