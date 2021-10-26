package dictionary

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

// Adapter Adapter适配器
type Adapter interface {
	GetDataSetByName(string) (*models.DataSet, error)
	GetSourcesByIds([]uint64) ([]*models.DataSource, error)
	GetMetricsByIds([]uint64) ([]*models.Metric, error)
	GetDimensionsByIds([]uint64) ([]*models.Dimension, error)
}

// AdapterOption Adapter配置
type AdapterOption struct {
	Type AdapterType
	Dsn  string
}

func NewAdapter(option *AdapterOption) (Adapter, error) {
	switch option.Type {
	case FILEAdapter:
		return newDictionaryAdapterByFile(option)
	default:
		return nil, fmt.Errorf("not supported adapter type %v", option.Type)
	}
}

// FileAdapter 文件适配器
type FileAdapter struct {
	Sets       []*models.DataSet    `yaml:"sets"       json:"sets"`
	Sources    []*models.DataSource `yaml:"sources"    json:"sources"`
	Metrics    []*models.Metric     `yaml:"metrics"    json:"metrics"`
	Dimensions []*models.Dimension  `yaml:"dimensions" json:"dimensions"`
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
	if err = adapter.AddId(); err != nil {
		return nil, err
	}
	if err = adapter.isValidAdapterCheck(); err != nil {
		return nil, err
	}
	return adapter, nil
}

func (d *FileAdapter) GetDataSetByName(name string) (*models.DataSet, error) {
	for _, data := range d.Sets {
		if data.Name == name {
			return data, nil
		}
	}
	return nil, fmt.Errorf("can not find '%v' data set", name)
}

func (d *FileAdapter) GetSourcesByIds(ids []uint64) ([]*models.DataSource, error) {
	idsMap := getIdsMap(ids)
	metricsSourcesIdsMap := make(map[uint64]bool)
	for _, metric := range d.Metrics {
		metricsSourcesIdsMap[metric.DataSourceID] = true
	}

	dimensionsSourcesIdsMap := make(map[uint64]bool)
	for _, dimension := range d.Dimensions {
		dimensionsSourcesIdsMap[dimension.DataSourceID] = true
	}

	result := make([]*models.DataSource, 0)
	for _, source := range d.Sources {
		_, ok := idsMap[source.ID]
		_, ok2 := metricsSourcesIdsMap[source.ID]
		_, ok3 := dimensionsSourcesIdsMap[source.ID]
		if ok && (ok2 || ok3) {
			result = append(result, source)
		}
	}
	return result, nil
}

func (d *FileAdapter) GetMetricsByIds(ids []uint64) ([]*models.Metric, error) {
	idsMap := getIdsMap(ids)
	metrics := make([]*models.Metric, 0)
	for _, metric := range d.Metrics {
		if _, ok := idsMap[metric.DataSourceID]; ok {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

// GetDimensionsByIds 通过ids筛选Dimensions信息
func (d *FileAdapter) GetDimensionsByIds(ids []uint64) ([]*models.Dimension, error) {
	idsMap := getIdsMap(ids)
	dimensions := make([]*models.Dimension, 0)
	for _, dimension := range d.Dimensions {
		if _, ok := idsMap[dimension.DataSourceID]; ok {
			dimensions = append(dimensions, dimension)
		}
	}
	return dimensions, nil
}

func (d *FileAdapter) isValidJoinOns(joinOns models.JoinOns) (id1, id2 uint64, err error) {
	in1, in2 := joinOns.ID()

	in1Map := getIdsMap(in1)
	in2Map := getIdsMap(in2)

	out1 := make(map[uint64]bool, 0)
	out2 := make(map[uint64]bool, 0)

	for _, dimension := range d.Dimensions {
		id := dimension.ID
		if _, ok := in1Map[id]; ok {
			out1[dimension.DataSourceID] = true
		}
		if _, ok := in2Map[id]; ok {
			out2[dimension.DataSourceID] = true
		}
	}

	if len(out1) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out1)
	}
	if len(out2) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out2)
	}

	for id := range out1 {
		id1 = id
	}

	for id := range out2 {
		id2 = id
	}
	return
}

func (d *FileAdapter) isValidSecondary(secondary *models.Secondary) error {
	id1, id2, err := d.isValidJoinOns(secondary.JoinOn)
	if err != nil {
		return err
	}
	if id1 != secondary.DataSourceID1 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id1, secondary.DataSourceID1)
	}
	if id2 != secondary.DataSourceID2 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id2, secondary.DataSourceID2)
	}
	return nil
}

func (d *FileAdapter) isValidDataSetSchema(schema *models.DataSetSchema) error {
	if _, err := schema.Tree(); err != nil {
		return err
	}

	for _, v := range schema.Secondary {
		if err := d.isValidSecondary(v); err != nil {
			return err
		}
	}
	return nil
}

// isValidDataSet 检查DataSet的合法性
func (d *FileAdapter) isValidDataSet(set *models.DataSet) error {
	return d.isValidDataSetSchema(set.Schema)
}

// isValidAdapterCheck 检查Adapter的合法性
func (d *FileAdapter) isValidAdapterCheck() error {
	for _, set := range d.Sets {
		if err := d.isValidDataSet(set); err != nil {
			return err
		}
	}
	return nil
}

func getIdsMap(ids []uint64) map[interface{}]interface{} {
	idsMap := make(map[interface{}]interface{})
	for _, id := range ids {
		idsMap[id] = true
	}
	return idsMap
}

func (d *FileAdapter) AddId() error {
	for i, source := range d.Sources {
		source.ID = uint64(i + 1)
	}
	for i, metric := range d.Metrics {
		metric.ID = uint64(i + 1)
	}
	for i, dimension := range d.Dimensions {
		dimension.ID = uint64(i + 1)
	}
	return nil
}