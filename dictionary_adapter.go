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
	BuildDataSetAdapter(string) (IAdapter, error)

	GetMetric() []*models.Metric

	GetDataSetByKey(string) (*models.DataSet, error)
	GetSourceByKey(string) (*models.DataSource, error)
	GetMetricByKey(string) (*models.Metric, error)
	GetDimensionByKey(string) (*models.Dimension, error)

	GetSourcesByKeys([]string) ([]*models.DataSource, error)
	GetMetricsByKeys([]string) ([]*models.Metric, error)
	GetMetricsBySourceKeys([]string) ([]*models.Metric, error)
	GetDimensionsByKeys([]string) ([]*models.Dimension, error)
	GetDimensionsBySourceKeys([]string) ([]*models.Dimension, error)
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

func (f *FileAdapter) BuildDataSetAdapter(key string) (IAdapter, error) {
	set, err := f.GetDataSetByKey(key)
	if err != nil {
		return nil, err
	}

	sKey := set.GetDataSource()
	sources, err := f.GetSourcesByKeys(sKey)
	if err != nil {
		return nil, err
	}

	metrics, err := f.GetMetricsBySourceKeys(sKey)
	if err != nil {
		return nil, err
	}

	dimensions, err := f.GetDimensionsBySourceKeys(sKey)
	if err != nil {
		return nil, err
	}

	adapter := &FileAdapter{
		Sets:       []*models.DataSet{set},
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
		if source.GetKey() == key {
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

func (f *FileAdapter) GetSourcesByKeys(key []string) ([]*models.DataSource, error) {
	set := getKeySet(key)
	var out []*models.DataSource
	for _, source := range f.Sources {
		if _, ok := set[source.GetKey()]; ok {
			out = append(out, source)
		}
	}
	return out, nil
}

func (f *FileAdapter) GetMetricsByKeys(key []string) ([]*models.Metric, error) {
	set := getKeySet(key)
	var out []*models.Metric
	for _, metric := range f.Metrics {
		if _, ok := set[metric.GetKey()]; ok {
			out = append(out, metric)
		}
	}
	return out, nil
}

func (f *FileAdapter) GetMetricsBySourceKeys(key []string) ([]*models.Metric, error) {
	set := getKeySet(key)
	var out []*models.Metric
	for _, metric := range f.Metrics {
		if _, ok := set[metric.DataSource]; ok {
			out = append(out, metric)
		}
	}
	return out, nil
}

func (f *FileAdapter) GetDimensionsByKeys(key []string) ([]*models.Dimension, error) {
	set := getKeySet(key)
	var out []*models.Dimension
	for _, dimension := range f.Dimensions {
		if _, ok := set[dimension.GetKey()]; ok {
			out = append(out, dimension)
		}
	}
	return out, nil
}

func (f *FileAdapter) GetDimensionsBySourceKeys(key []string) ([]*models.Dimension, error) {
	set := getKeySet(key)
	var out []*models.Dimension
	for _, dimension := range f.Dimensions {
		if _, ok := set[dimension.DataSource]; ok {
			out = append(out, dimension)
		}
	}
	return out, nil
}

func (f *FileAdapter) isValidJoin(join *models.DataSetJoin) error {
	var key1, key2 []string
	for _, on := range join.JoinOn {
		k1 := fmt.Sprintf("%v.%v", join.DataSource1, on.Dimension1)
		k2 := fmt.Sprintf("%v.%v", join.DataSource2, on.Dimension2)
		key1 = append(key1, k1)
		key2 = append(key2, k2)
	}

	d1, _ := f.GetDimensionsByKeys(key1)
	if len(d1) != len(join.JoinOn) {
		return fmt.Errorf("invalid dataset join setting=%v, found dimension=%v", key1, d1)
	}
	d2, _ := f.GetDimensionsByKeys(key2)
	if len(d2) != len(join.JoinOn) {
		return fmt.Errorf("invalid dataset join setting=%v, found dimension=%v", key2, d2)
	}

	source, _ := f.GetSourcesByKeys([]string{join.DataSource1, join.DataSource2})
	if len(source) != 2 {
		return fmt.Errorf("invalid dataset join setting, found source=%v", source)
	}

	return nil
}

func (f *FileAdapter) isValidJoinTopologyGraph(set *models.DataSet) error {
	_, err := set.JoinTopologyGraph()
	return err
}

func (f *FileAdapter) isValidDataSet(set *models.DataSet) error {
	if err := f.isValidJoinTopologyGraph(set); err != nil {
		return err
	}
	for _, join := range set.Join {
		if err := f.isValidJoin(join); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileAdapter) isValidMetric(metric *models.Metric) error {
	// TODO
	return nil
}

func (f *FileAdapter) isValidDimension(dimension *models.Dimension) error {
	// TODO
	return nil
}

func (f *FileAdapter) isValid() error {
	for _, set := range f.Sets {
		if err := f.isValidDataSet(set); err != nil {
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
