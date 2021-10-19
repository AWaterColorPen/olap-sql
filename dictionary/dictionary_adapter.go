package dictionary

import (
	"fmt"
	"io/ioutil"

	"github.com/awatercolorpen/olap-sql/api/models"
	"gopkg.in/yaml.v2"
)

type AdapterType string

const (
	DBadapter   AdapterType = "DB"
	FILEadapter AdapterType = "FILE"
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
	// 根据不同的Type去实例化不同的Adapter
	switch option.Type {
	case DBadapter:
		// TODO
		return newDictionaryAdapterByDB(option)
	case FILEadapter:
		return newDictionaryAdapterByYaml(option)
	default:
		return nil, fmt.Errorf("adapter type error")
	}
}

// FileAdapter 文件适配器
type FileAdapter struct {
	// TODO
	Sets       []*models.DataSet    `yaml:"sets"`
	Sources    []*models.DataSource `yaml:"sources"`
	Metrics    []*models.Metric     `yaml:"metrics"`
	Dimensions []*models.Dimension  `yaml:"dimensions""`
}

func (d *FileAdapter) move() {
	fmt.Println("qwq")
}
func (d *FileAdapter) Create(item interface{}) error {
	switch v := item.(type) {
	case *models.DataSet:
		if err := d.isValidDataSetSchema(v.Schema); err != nil {
			return err
		}
		d.Sets = append(d.Sets, v)
	case []*models.DataSet:
		for _, i := range item.([]*models.DataSet) {
			if err := d.isValidDataSetSchema(i.Schema); err != nil {
				return err
			}
			d.Sets = append(d.Sets, i)
		}
	case *models.DataSource:
		d.Sources = append(d.Sources, v)
	case []*models.DataSource:
		d.Sources = append(d.Sources, v...)
	case *models.Metric:
		d.Metrics = append(d.Metrics, v)
	case []*models.Metric:
		d.Metrics = append(d.Metrics, v...)
	case *models.Dimension:
		d.Dimensions = append(d.Dimensions, v)
	case []*models.Dimension:
		d.Dimensions = append(d.Dimensions, v...)
	}
	return nil
}

func newDictionaryAdapterByDB(option *AdapterOption) (Adapter, error) {
	return nil, fmt.Errorf("DB type unsupport now")
}

func newDictionaryAdapterByYaml(option *AdapterOption) (*FileAdapter, error) {
	adapter := &FileAdapter{}
	yamlFile, err := ioutil.ReadFile(option.Dsn)
	if err != nil {
		return nil, fmt.Errorf("file read error")
	}
	if err := yaml.Unmarshal(yamlFile, adapter); err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed")
	}
	return adapter, nil
}

func (d *FileAdapter) GetDataSetByName(name string) (*models.DataSet, error) {
	for _, data := range d.Sets {
		if data.Name == name {
			return checkDataSetActive(data)
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

func checkDataSetActive(set *models.DataSet) (*models.DataSet, error) {
	if set.Schema == nil {
		return nil, fmt.Errorf("schema is nil for data_set %v", set.Name)
	}
	return set, nil
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
